# Affiliate

## Event checking and notifying on creating trading order

<img src='https://mm.meta.etop.vn/g?
sequenceDiagram;
participant U as User;
participant TAPI as TradingAPI;
participant OPM as OrderTradingPM;
participant OLOGIC as OrderLogic;
participant AG as AffiliateAggregate;
U ->> TAPI: TradingCreateOrder;
TAPI ->> OPM: TradingOrderCreatingEvent;
OPM ->> AG: TradingOrderCreatingCommand;
alt success;
AG -->> OPM: OK;
OPM -->> TAPI: OK;
TAPI ->> OLOGIC: CreateOrder;
OLOGIC -->> TAPI: newOrder;
TAPI ->> OPM: TradingOrderCreatedEvent{orderID, referralCode};
OPM ->> AG: OnTradingOrderCreatedCommand{orderID, referralCode};
AG -->> OPM: OK;
OPM -->> TAPI: OK;
TAPI -->> U: newOrder;
else failure;
AG -->> OPM: TradingOrderInvalidError;
OPM -->> TAPI: TradingOrderInvalidError;
TAPI -->> U: CreateOrderFailed;
end;
'>