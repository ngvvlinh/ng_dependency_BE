#!/usr/bin/env bash
set -e

: ${1?"expect PREFIX (should be d for local)"}

PREFIX=${1}_

# Run this file inside the kafka container

kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions   8 --topic ${PREFIX}pgrid_account
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions   8 --topic ${PREFIX}pgrid_address
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions   8 --topic ${PREFIX}pgrid_etop_category
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions  64 --topic ${PREFIX}pgrid_fulfillment
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions  64 --topic ${PREFIX}pgrid_order_external
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions  64 --topic ${PREFIX}pgrid_order_source
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions  64 --topic ${PREFIX}pgrid_order
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions   8 --topic ${PREFIX}pgrid_product_brand
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions  64 --topic ${PREFIX}pgrid_product_external
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions   8 --topic ${PREFIX}pgrid_product_source_category_external
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions   8 --topic ${PREFIX}pgrid_product_source_category
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions   8 --topic ${PREFIX}pgrid_product_source
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions  64 --topic ${PREFIX}pgrid_product
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions   8 --topic ${PREFIX}pgrid_shop_collection
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions  64 --topic ${PREFIX}pgrid_shop_product
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions   8 --topic ${PREFIX}pgrid_shop
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions   8 --topic ${PREFIX}pgrid_supplier
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions   8 --topic ${PREFIX}pgrid_user
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions  64 --topic ${PREFIX}pgrid_variant_external
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions  64 --topic ${PREFIX}pgrid_variant
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions   8 --topic ${PREFIX}pgrid_money_transaction_shipping
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions   64 --topic ${PREFIX}pgrid_notification

# UpdateInfo 2018-06-12

# Internal control channel
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions   1 --topic ${PREFIX}intctl
