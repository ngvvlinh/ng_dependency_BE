#!/usr/bin/env bash
set -e

: ${1?"expect PREFIX (should be d for local)"}

PREFIX=${1}_

# Run this file inside the kafka container

kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions  64 --topic ${PREFIX}pgrid_fulfillment
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions  64 --topic ${PREFIX}pgrid_order
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions  8 --topic ${PREFIX}pgrid_money_transaction_shipping
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions   64 --topic ${PREFIX}pgrid_notification
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions 64 --topic ${PREFIX}pgrid_shop_product

# UpdateInfo 2018-06-12

# Internal control channel
kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions   1 --topic ${PREFIX}intctl
