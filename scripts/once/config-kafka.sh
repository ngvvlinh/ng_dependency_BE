#!/usr/bin/env bash

: ${1?"expect PREFIX (should be d for local)"}

PREFIX=${1}_

create_topic() {
    : ${1?"Usage: create_topic PARTITIONS TOPIC_NAME"}
    : ${2?"Usage: create_topic PARTITIONS TOPIC_NAME"}

    kafka-topics --create --zookeeper zookeeper:2181 --replication-factor 1 --partitions "$1" --topic "${PREFIX}$2"
    kafka-topics --alter  --zookeeper zookeeper:2181                        --partitions "$1" --topic "${PREFIX}$2"
}

# Run this file inside the kafka container

create_topic 64 pgrid_fulfillment
create_topic 64 pgrid_order
create_topic 64 pgrid_order
create_topic 8  pgrid_money_transaction_shipping
create_topic 8  pgrid_money_transaction_shipping
create_topic 64 pgrid_notification
create_topic 64 pgrid_notification
create_topic 64 pgrid_shop_product
create_topic 64 pgrid_shop_product
create_topic 64 pgrid_shop_collection
create_topic 64 pgrid_shop_collection
create_topic 64 pgrid_shop_product_collection
create_topic 64 pgrid_shop_product_collection
create_topic 64 pgrid_shop_variant
create_topic 64 pgrid_shop_variant
create_topic 64 pgrid_inventory_variant
create_topic 64 pgrid_inventory_variant
create_topic 64 pgrid_shop_trader_address
create_topic 64 pgrid_shop_trader_address
create_topic 64 pgrid_shop_customer
create_topic 64 pgrid_shop_customer
create_topic 64 pgrid_shop_customer_group
create_topic 64 pgrid_shop_customer_group
create_topic 64 pgrid_shop_customer_group_customer
create_topic 64 pgrid_shop_customer_group_customer
create_topic 64 pgrid_fb_external_user
create_topic 64 pgrid_fb_external_comment
create_topic 64 pgrid_fb_external_conversation
create_topic 64 pgrid_fb_customer_conversation
create_topic 64 pgrid_fb_external_message
create_topic 64 pgrid_fb_external_user_fabo
create_topic 64 pgrid_fb_external_comment_fabo
create_topic 64 pgrid_fb_external_conversation_fabo
create_topic 64 pgrid_fb_external_message_fabo
create_topic 64 pgrid_fb_customer_conversation_fabo
create_topic 64 pgrid_shipnow_fulfillment


# UpdateInfo 2018-06-12

# Internal control channel
create_topic  1 intctl
