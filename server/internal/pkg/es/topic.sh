curl -XPUT 127.0.0.1:9200/bbsgo_topic_1 -H 'Content-Type: application/json' -d'
{
    "mappings": {
        "properties": {
            "id": {
                "type": "long",
                "store": true
            },
            "nodeId": {
                "type": "long",
                "store": true
            },
            "userId": {
                "type": "long",
                "store": true
            },
            "nickname": {
                "type": "text",
                "analyzer": "ik_max_word",
                "search_analyzer": "ik_smart"
            },
            "title": {
                "type": "text",
                "analyzer": "ik_max_word",
                "search_analyzer": "ik_smart"
            },
            "content": {
                "type": "text",
                "analyzer": "ik_max_word",
                "search_analyzer": "ik_smart"
            },
            "tags": {
                "type": "text",
                "analyzer": "ik_max_word",
                "search_analyzer": "ik_smart",
                "fields": {
                    "keyword": {
                        "type": "keyword",
                        "ignore_above": 256
                    }
                }
            },
            "recommend": {
                "type": "boolean"
            },
            "status": {
                "type": "integer"
            },
            "createTime": {
                "type": "long",
                "store": true
            }
        }
    }
}'
