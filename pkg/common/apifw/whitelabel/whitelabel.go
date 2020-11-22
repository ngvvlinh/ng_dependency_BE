package whitelabel

import (
	"context"
	"strings"

	"o.o/api/main/identity"
	"o.o/backend/pkg/common/headers"
	"o.o/capi/dot"
	"o.o/common/l"
	"o.o/common/xerrors"
)

var ll = l.New()

type whitelabelKey struct{}

type WhiteLabel struct {
	etop          *WL
	partners      []*WL
	partnersByID  map[dot.ID]*WL
	partnersByKey map[string]*WL
}

func New(drivers []*WL) *WhiteLabel {
	w := &WhiteLabel{}
	w.partners = drivers
	w.partnersByID = make(map[dot.ID]*WL)
	w.partnersByKey = make(map[string]*WL)
	for _, driver := range drivers {
		w.partnersByKey[driver.Key] = driver
		w.partnersByID[driver.ID] = driver
	}
	w.etop = w.partnersByID[0]
	return w
}

func (w *WhiteLabel) VerifyPartners(ctx context.Context, identityQuery identity.QueryBus) error {
	query := &identity.ListPartnersForWhiteLabelQuery{}
	if err := identityQuery.Dispatch(ctx, query); err != nil {
		return err
	}

	dbKeys := make(map[string]bool)
	for _, partner := range query.Result {
		p := w.partnersByID[partner.ID]
		if p == nil {
			return xerrors.Errorf(xerrors.Internal, nil, "white label partner %v not found", partner.ID)
		}
		if p.Key != partner.WhiteLabelKey {
			return xerrors.Errorf(xerrors.Internal, nil, "white label key of partner %v not match (partner.white_label_key=%v driver.key=%v)", partner.ID, partner.WhiteLabelKey, p.Key)
		}
		dbKeys[partner.WhiteLabelKey] = true
	}

	for _, p := range w.partners[1:] { // ignore the first one - etop
		if !dbKeys[p.Key] {
			ll.S.Errorf("white label key %v not found in db", p.Key)
		}
	}
	return nil
}

func (w *WhiteLabel) ByPartnerID(id dot.ID) *WL {
	return w.byPartnerID(id).Clone()
}

func (w *WhiteLabel) byPartnerID(id dot.ID) *WL {
	p := w.partnersByID[id]
	if p == nil {
		return w.etop.Clone()
	}
	return p.Clone()
}

func (w *WhiteLabel) ByPartnerKey(key string) *WL {
	p := w.partnersByKey[key]
	if p == nil {
		return w.etop.Clone()
	}
	return p.Clone()
}

func (w *WhiteLabel) ByContext(ctx context.Context) *WL {
	wl := ctx.Value(whitelabelKey{})
	if wl == nil {
		return w.etop.Clone()

		// MUSTDO(vu): enable this
		// ll.Panic("whitelabel context should not be called here")
	}
	return wl.(*WL)
}

func (w *WhiteLabel) WrapContext(ctx context.Context, partnerID dot.ID) context.Context {
	return context.WithValue(ctx, whitelabelKey{},
		w.fromContext(ctx, partnerID).Clone())
}

func (w *WhiteLabel) WrapContextByPartnerID(ctx context.Context, partnerID dot.ID) context.Context {
	return context.WithValue(ctx, whitelabelKey{}, w.fromPartnerID(partnerID).Clone())
}

func (w *WhiteLabel) fromPartnerID(partnerID dot.ID) *WL {
	return w.ByPartnerID(partnerID)
}

func (w *WhiteLabel) fromContext(ctx context.Context, partnerID dot.ID) *WL {
	if partnerID != 0 {
		return w.byPartnerID(partnerID)
	}
	header := headers.GetHeader(ctx)
	if header == nil {
		panic("no http header")
	}
	host := header.Get("X-Forwarded-Host")
	return w.fromHost(host)
}

func (w *WhiteLabel) fromHost(host string) *WL {
	parts := strings.SplitN(host, ".", 2)
	for _, p := range w.partners {
		if p.IgnoreParseFromHost {
			continue
		}
		// itopx.vn
		if p.Host == host {
			return p
		}
		if isWhitelabelKey(p.Key, parts[0]) {
			return p
		}
	}
	return w.etop
}

func isWhitelabelKey(key, subdomain string) bool {
	// itopx.d.etop.vn
	if subdomain == key {
		return true
	}
	// itopx-next.etop.vn
	return len(subdomain) > len(key) &&
		subdomain[0:len(key)] == key &&
		subdomain[len(key)] == '-'
}
