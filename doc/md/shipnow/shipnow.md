# Shipnow

## Confirm Shipnow

<img src='https://mm.meta.etop.vn/g?
sequenceDiagram;
participant U as User;
participant S as ShopAPI;
participant SA as ShipnowAggregate;
participant CPM as ConfirmationPM;
participant OA as OrderAggregate;
participant CA as CarrierAggregate;
U ->> S: ConfirmShipnow;
S ->> SA: ConfirmShipnowCommand;
SA -->> CPM: ConfirmationRequestedEvent;
CPM ->> OA: PrepareOrdersForFulfillmentCommand;
OA -->> CPM: OK;
CPM ->> +CA: CreateExternalShipmentCommand;
alt success;
CPM ->> OA: CommitOrdersForFulfillmentCommand;
CA -->> CPM: ExternalShipmentCreatedEvent;
CPM ->> SA: CommitShipnowCreationCommand;
else failure;
CA -->> CPM: ExternalShipmentCreationRejectedEvent;
CPM ->> OA: ReleaseOrdersForFulfillmentCommand;
CPM ->> SA: CancelShipnowCreationCommand;
end;
SA -->> S: change/shipnow {id=123};
'>
