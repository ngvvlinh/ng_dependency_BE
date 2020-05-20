package webserver

import (
	"bytes"
	"context"
	"html/template"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"o.o/api/main/address"
	"o.o/api/main/catalog"
	"o.o/api/main/location"
	"o.o/api/top/int/types"
	"o.o/api/top/types/etc/gender"
	"o.o/api/top/types/etc/order_source"
	"o.o/api/top/types/etc/payment_method"
	"o.o/api/webserver"
	"o.o/backend/com/main/identity/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/authorize/claims"
	logicorder "o.o/backend/pkg/etop/logic/orders"
	"o.o/capi/dot"
)

type SessionCart struct {
	Products    []*CartProductWithOneVariant `json:"products"`
	TotalAmount int                          `json:"total_amount"`
	TotalCount  int                          `json:"total_count"`
}

type CartProductWithOneVariant struct {
	Product *catalog.ShopProductWithVariants `json:"product"`
	Variant *catalog.ShopVariant             `json:"variant"`
	Count   int                              `json:"count"`
}

type CartResponse struct {
	Cart     *SessionCart `json:"cart"`
	DecsHTML string       `json:"decs_html"`
}

type UpdateCartReqest struct {
	Cart []*UpdateCartVariant `json:"cart"`
}

type UpdateCartVariant struct {
	VairantID string `json:"variant_id"`
	Quantity  string `json:"quantity"`
}

var prefixSessionOrder = "order"

func (s *Server) CartTotalCount(c echo.Context) error {
	sess, err := getSessionCart(c)
	if err != nil {
		return err
	}
	return c.String(200, strconv.Itoa(sess.TotalCount))
}

func (s *Server) CreateOrder(c echo.Context) error {
	fullname := c.FormValue("full_name")
	phone := c.FormValue("phone")
	province := c.FormValue("province_code")
	district := c.FormValue("district_code")
	ward := c.FormValue("ward_code")
	address := c.FormValue("address1")
	note := c.FormValue("note")
	if fullname == "" || phone == "" || province == "" || district == "" || address == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "thiếu thông tin tạo đơn hàng")
	}
	shopID, err := GetShopID(c)
	if err != nil {
		return err
	}
	sessionCart, err := getSessionCart(c)
	if err != nil {
		return err
	}
	// check remove product variant
	variantIDs := []dot.ID{}
	for _, v := range sessionCart.Products {
		variantIDs = append(variantIDs, v.Variant.VariantID)
	}
	if len(variantIDs) > 0 {
		queryVariant := &catalog.ListShopVariantsByIDsQuery{
			IDs:    variantIDs,
			ShopID: sessionCart.Products[0].Variant.ShopID,
		}
		err = catelogQueryBus.Dispatch(c.Request().Context(), queryVariant)
		if err != nil {
			return err
		}
		var mapVariants = make(map[dot.ID]*catalog.ShopVariant)
		for _, v := range queryVariant.Result.Variants {
			mapVariants[v.VariantID] = v
		}
		var listVariantRemoved []*CartProductWithOneVariant
		for _, v := range sessionCart.Products {
			if mapVariants[v.Variant.VariantID] == nil {
				listVariantRemoved = append(listVariantRemoved, v)
			}
		}
		if len(listVariantRemoved) > 0 {
			return c.JSON(202, listVariantRemoved)
		}
	}

	queryLocation := &location.FindOrGetLocationQuery{
		ProvinceCode: province,
		DistrictCode: district,
	}
	if ward != "undefined" && ward != "" {
		queryLocation.WardCode = ward
	}
	err = locationBus.Dispatch(c.Request().Context(), queryLocation)
	if err != nil {
		return err
	}
	claimsShop := claims.ShopClaim{
		UserClaim:          claims.UserClaim{},
		CommonAccountClaim: claims.CommonAccountClaim{},
		Actions:            nil,
		Shop:               &model.Shop{ID: shopID},
	}
	lines := []*types.CreateOrderLine{}
	totalCount := 0
	totalAmount := 0
	for _, v := range sessionCart.Products {
		totalCount += v.Count
		totalAmount += v.Variant.RetailPrice * v.Count
		lines = append(lines, &types.CreateOrderLine{
			VariantId:    v.Variant.VariantID,
			ProductName:  v.Product.Name,
			Quantity:     v.Count,
			ListPrice:    v.Variant.ListPrice,
			RetailPrice:  v.Variant.RetailPrice,
			PaymentPrice: v.Variant.RetailPrice,
			Attributes:   v.Variant.Attributes,
		})
	}

	customerAddress := &types.OrderAddress{
		ExportedFields: nil,
		FullName:       fullname,
		Phone:          phone,
		Province:       queryLocation.Result.Province.Name,
		District:       queryLocation.Result.District.Name,

		Address1:     address,
		ProvinceCode: province,
		DistrictCode: district,
	}
	if queryLocation.Result.Ward != nil {
		customerAddress.Ward = queryLocation.Result.Ward.Name
		customerAddress.WardCode = ward
	}
	cmdCreateOrder := &types.CreateOrderRequest{
		OrderNote:     note,
		Lines:         lines,
		Source:        order_source.Ecomify,
		PaymentMethod: payment_method.COD,
		Customer: &types.OrderCustomer{
			FullName: fullname,
			Phone:    phone,
			Gender:   gender.Other,
		},
		CustomerAddress: customerAddress,
		TotalItems:      totalCount,
		TotalAmount:     totalAmount,
		BasketValue:     totalAmount,
		PreOrder:        true,
	}

	order, err := logicorder.CreateOrder(c.Request().Context(), &claimsShop, nil, cmdCreateOrder, nil, 0)
	if err != nil {
		return err
	}
	err = clearSessionCart(c)
	if err != nil {
		return err
	}
	err = setSessionOrder(c, order.Code)
	if err != nil {
		return err
	}
	return c.JSON(200, &CartResponse{Cart: sessionCart})
}

func (s *Server) CartRemoveVariant(c echo.Context) error {
	variantIDString := c.FormValue("variant_id")
	variantID, err := dot.ParseID(variantIDString)
	if err != nil {
		return err
	}
	sessionCart, err := getSessionCart(c)
	if err != nil {
		return err
	}
	sessionCart.TotalAmount = 0
	sessionCart.TotalCount = 0
	var products []*CartProductWithOneVariant
	for _, product := range sessionCart.Products {
		if product.Variant.VariantID != variantID {
			sessionCart.TotalCount += product.Count
			sessionCart.TotalAmount += product.Count * product.Variant.RetailPrice
			products = append(products, product)
		}
	}
	sessionCart.Products = products
	err = setSessionCart(c, sessionCart)
	if err != nil {
		return err
	}
	tmpl, err := template.New("cart").Parse(cartHTMLTCart)
	if err != nil {
		return err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, &CartResponse{Cart: sessionCart})
	if err != nil {
		return err
	}

	return c.HTML(200, tpl.String())
}

func (s *Server) CartQuickAddProduct(c echo.Context) error {
	shopID, err := GetShopID(c)
	if err != nil {
		return nil
	}
	productIDString := c.FormValue("product_id")
	productID, err := dot.ParseID(productIDString)
	if err != nil {
		return err
	}

	query := &webserver.GetWsProductByIDQuery{
		ShopID: shopID,
		ID:     productID,
	}
	err = webserverQueryBus.Dispatch(c.Request().Context(), query)
	if err != nil {
		return err
	}

	if !query.Result.Appear {
		c.Response().Status = http.StatusNotFound
		return nil
	}
	isAdd := false
	var shopVariant = query.Result.Product.Variants[0]
	if len(query.Result.Product.Variants) < 1 {
		c.Response().Status = http.StatusNotFound
		return nil
	}

	sessionCart, err := getSessionCart(c)
	if err != nil {
		return err
	}
	sessionCart.TotalCount++
	variantID := query.Result.Product.Variants[0].VariantID
	for k, v := range sessionCart.Products {
		if v.Variant.VariantID == variantID {
			sessionCart.Products[k].Count++
			sessionCart.TotalAmount += v.Variant.RetailPrice
			isAdd = true
		}
	}
	if !isAdd {
		sessionCart.Products = append(sessionCart.Products, &CartProductWithOneVariant{
			Product: query.Result.Product,
			Variant: shopVariant,
			Count:   1,
		})
		sessionCart.TotalAmount += shopVariant.RetailPrice
	}
	err = setSessionCart(c, sessionCart)
	if err != nil {
		return err
	}
	tmpl, err := template.New("cart").Parse(cartHTMLTCart)
	if err != nil {
		return err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, &CartResponse{Cart: sessionCart})
	if err != nil {
		return err
	}

	return c.HTML(200, tpl.String())
}

func (s *Server) CartAddProduct(c echo.Context) error {
	shopID, err := GetShopID(c)
	if err != nil {
		return nil
	}
	productIDString := c.FormValue("product_id")
	variantIDString := c.FormValue("variant_id")
	countString := c.FormValue("count")

	productID, err := dot.ParseID(productIDString)
	if err != nil {
		return err
	}

	variantID, err := dot.ParseID(variantIDString)
	if err != nil {
		return err
	}

	count, err := strconv.Atoi(countString)
	if err != nil {
		return err
	}

	query := &webserver.GetWsProductByIDQuery{
		ShopID: shopID,
		ID:     productID,
	}
	err = webserverQueryBus.Dispatch(c.Request().Context(), query)
	if err != nil {
		return err
	}
	sessionCart, err := getSessionCart(c)
	if err != nil {
		return err
	}
	isAdd := false
	var shopVariant *catalog.ShopVariant
	for _, v := range query.Result.Product.Variants {
		if variantID == v.VariantID {
			shopVariant = v
		}
	}
	for k, v := range sessionCart.Products {
		if v.Variant.VariantID == variantID {
			sessionCart.Products[k].Count = sessionCart.Products[k].Count + count
			sessionCart.TotalCount += count
			sessionCart.TotalAmount += count * v.Variant.RetailPrice
			isAdd = true
			break
		}
	}
	if !isAdd {
		var cartAddProduct = &CartProductWithOneVariant{
			Product: query.Result.Product,
			Variant: shopVariant,
			Count:   count,
		}
		sessionCart.Products = append(sessionCart.Products, cartAddProduct)
		sessionCart.TotalAmount += cartAddProduct.Count * cartAddProduct.Variant.RetailPrice
		sessionCart.TotalCount += cartAddProduct.Count
	}

	err = setSessionCart(c, sessionCart)
	if err != nil {
		return err
	}

	tmpl, err := template.New("cart").Parse(cartHTMLTCart)
	if err != nil {
		return err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, &CartResponse{Cart: sessionCart})
	if err != nil {
		return err
	}

	return c.HTML(200, tpl.String())
}

func (s *Server) CartUpdateAllListProduct(c echo.Context) error {
	sessionCart, err := getSessionCart(c)
	if err != nil {
		return err
	}
	u := new(UpdateCartReqest)
	if err := c.Bind(u); err != nil {
		return err
	}
	var sessionCartUpdate = &SessionCart{}
	totalCount := 0
	totalAmount := 0
	for _, v := range u.Cart {
		variantID, err := dot.ParseID(v.VairantID)
		if err != nil {
			return err
		}
		quantity, err := strconv.Atoi(v.Quantity)
		if err != nil {
			return err
		}
		for _, cartItem := range sessionCart.Products {
			if cartItem.Variant.VariantID == variantID {
				cartItem.Count = quantity
				totalCount += quantity
				totalAmount += quantity * cartItem.Variant.RetailPrice
				sessionCartUpdate.Products = append(sessionCartUpdate.Products, cartItem)
				break
			}
		}
	}
	sessionCartUpdate.TotalCount = totalCount
	sessionCartUpdate.TotalAmount = totalAmount
	err = setSessionCart(c, sessionCartUpdate)
	if err != nil {
		return err
	}

	tmpl, err := template.New("cart").Parse(cartHTMLTCart)
	if err != nil {
		return err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, &CartResponse{Cart: sessionCartUpdate})
	if err != nil {
		return err
	}
	return c.HTML(200, tpl.String())
}

func (s *Server) Cart(c echo.Context) error {
	indexData, err := GetIndexData(c)
	if err != nil {
		return err
	}
	err = c.Render(http.StatusOK, "header.html", indexData)
	if err != nil {
		return err
	}
	if len(indexData.Cart.Products) == 0 {
		err = c.Render(http.StatusOK, "cart-empty.html", indexData)
		if err != nil {
			return err
		}
	} else {
		err = c.Render(http.StatusOK, "cart.html", indexData)
		if err != nil {
			return err
		}
	}
	err = c.Render(http.StatusOK, "footer.html", indexData)
	if err != nil {
		return err
	}
	return nil
}

func getSessionOrder(c echo.Context) (string, error) {
	redisKey, err := GetCookieKey(c)
	if err != nil {
		return "", err
	}
	var sessionOrder = ""
	err = redisStore.Get(redisKey+prefixSessionOrder, &sessionOrder)
	if err != nil {
		return "", err
	}
	return sessionOrder, nil
}

func setSessionOrder(c echo.Context, orderCode string) error {
	redisKey, err := GetCookieKey(c)
	if err != nil {
		return err
	}
	err = redisStore.SetWithTTL(redisKey+prefixSessionOrder, orderCode, 30*24*60*60)
	if err != nil {
		return err
	}
	return nil
}

func getSessionCart(c echo.Context) (*SessionCart, error) {
	redisKey, err := GetCookieKey(c)
	if err != nil {
		return nil, err
	}
	var sessionCart = &SessionCart{}
	err = redisStore.Get(redisKey, sessionCart)
	if err != nil && err == redis.ErrNil {
		return &SessionCart{}, nil
	}
	if err != nil {
		return nil, err
	}
	return sessionCart, nil
}

func setSessionCart(c echo.Context, cart *SessionCart) error {
	redisKey, err := GetCookieKey(c)
	if err != nil {
		return err
	}
	err = redisStore.SetWithTTL(redisKey, cart, 30*24*60*60)
	if err != nil {
		return err
	}
	return nil
}

func clearSessionCart(c echo.Context) error {
	redisKey, err := GetCookieKey(c)
	if err != nil {
		return err
	}
	err = redisStore.Del(redisKey)
	if err != nil && err != redis.ErrNil {
		return err
	}
	return nil
}

func (s *Server) getAdressValue(ctx context.Context, arg *address.Address) error {

	if arg == nil {
		return nil
	}
	if arg.Province == "" || arg.ProvinceCode == "" {
		return cm.Error(cm.InvalidArgument, "Missing province information", nil)
	}
	if arg.District == "" || arg.DistrictCode == "" {
		return cm.Error(cm.InvalidArgument, "Missing district information", nil)
	}
	query := &location.FindOrGetLocationQuery{
		Province:     arg.Province,
		District:     arg.District,
		ProvinceCode: arg.ProvinceCode,
		DistrictCode: arg.DistrictCode,
	}
	if arg.WardCode != "" {
		query.WardCode = arg.WardCode
	}
	if err := locationBus.Dispatch(ctx, query); err != nil {
		return err
	}
	return nil
}

var cartHTMLTCart = `
{{if .Cart}}
<div class="dropcart__products-list">
	{{range .Cart.Products}}
		<div class="dropcart__product">
			<div class="dropcart__product-image">
				<a href="/product/product-{{.Product.ProductID}}"><img src="{{ $length := len .Product.ImageURLs }}
		{{ if eq $length 0 }}
			https://shop.d.etop.vn/assets/images/placeholder_medium.png
		{{else}}
			{{index .Product.ImageURLs 0}}
		{{end}}" alt=""> </a>
			</div>
			<div class="dropcart__product-info">
				<div class="dropcart__product-name"><a href="/product/product-{{.Product.ProductID}}">{{.Product.Name}}</a></div>
				{{if .Variant.Attributes }}
					<ul class="dropcart__product-options">
					   {{range .Variant.Attributes}}
						<li>{{.Name}}: {{.Value}}</li>
						{{end}}
					</ul>
				{{end}}
				<div class="dropcart__product-meta">
					<span class="dropcart__product-quantity">{{.Count}}</span> ×
					<span class="dropcart__product-price">{{.Variant.RetailPrice}}đ</span>
				</div>
			</div>
			<button type="button" class="dropcart__product-remove btn btn-light btn-sm btn-svg-icon" onclick="removerCart({{.Variant.VariantID}})">
				<i class="fa fa-times" aria-hidden="true" ></i>
			</button>
		</div>
	{{end}}
</div>
{{end}}
<div class="dropcart__totals">
<table>
	<tr>
		<th>Tổng cộng</th>
		{{if eq .Cart nil}}
			<td>0</td>
		{{else}}
			<td>{{.Cart.TotalAmount}}</td>
		{{end}}
	</tr>
</table>
</div>
<div class="dropcart__buttons">
<a class="btn btn-secondary px-0" href="/cart">Xem giỏ hàng</a>
<a class="btn btn-primary" href="/checkout">Đặt hàng</a>
</div>
`
