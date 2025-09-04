# MeiliSearch Migration Guide

This guide covers the migration from Bleve to MeiliSearch for search functionality in bbs-go.

## Overview

The application now supports both Bleve and MeiliSearch search engines with backward compatibility. MeiliSearch offers better performance, scalability, and search capabilities compared to Bleve.

## Prerequisites

1. Install and run MeiliSearch:
   ```bash
   # Option 1: Using Docker
   docker run -it --rm -p 7700:7700 -v meili_data:/meili_data getmeili/meilisearch:latest

   # Option 2: Download binary from https://github.com/meilisearch/meilisearch/releases
   ./meilisearch --master-key "123456" --db-path ./meili_data
   ```

2. MeiliSearch should be running on `localhost:7700` with master key `123456`

## Configuration

Update your `bbs-go.yaml` configuration file to include MeiliSearch settings:

```yaml
# MeiliSearch configuration
MeiliSearch:
  Host: localhost         # MeiliSearch server host
  Port: 7700             # MeiliSearch server port
  APIKey: "123456"       # MeiliSearch master key
  Index: topics          # Index name for topics
  Enabled: true          # Enable MeiliSearch (set to false to use Bleve instead)
```

### Configuration Options

- `Host`: MeiliSearch server hostname (default: localhost)
- `Port`: MeiliSearch server port (default: 7700)  
- `APIKey`: MeiliSearch master key for authentication (default: 123456)
- `Index`: Name of the index to store topics (default: topics)
- `Enabled`: Boolean flag to enable/disable MeiliSearch (default: true)

## Features

### Search Capabilities

MeiliSearch provides the following search features:

- **Full-text search** across title, content, nickname, and tags
- **Filtering** by node ID, user ID, status, and recommendation status
- **Time-based filtering** (last day, week, month, year)
- **Highlighting** of search terms in results
- **Pagination** support
- **Real-time indexing** of new/updated topics

### API Endpoints

1. **Search Topics**: `/api/search/topic`
   - Parameters:
     - `keyword`: Search term (optional)
     - `nodeId`: Filter by node ID (optional, -1 for recommended topics)
     - `timeRange`: Time filter (0=all, 1=day, 2=week, 3=month, 4=year)
     - `cursor`: Page number for pagination

2. **Reindex Topics**: `/api/search/reindex`
   - Rebuilds the search index from existing topics

## Migration Process

1. **Start MeiliSearch server** with the correct configuration
2. **Update configuration** to enable MeiliSearch
3. **Restart the application** to initialize MeiliSearch
4. **Reindex existing topics** by calling `/api/search/reindex`
5. **Test search functionality** using `/api/search/topic`

## Code Architecture

### Clean and Modular Design

The implementation follows clean architecture principles:

- **`meilisearch.go`**: Contains all MeiliSearch-specific logic
- **`search.go`**: Provides a unified interface that automatically routes to MeiliSearch or Bleve based on configuration
- **`common.go`**: Shared data structures and helper functions

### Key Components

1. **MeiliSearchClient**: Manages connection and operations with MeiliSearch
2. **Unified Interface**: Transparent switching between search engines
3. **Document Mapping**: Converts internal topic structure to search documents
4. **Batch Processing**: Efficient reindexing with configurable batch sizes

## Testing

Run the included test script to verify MeiliSearch integration:

```bash
go run test_meilisearch.go
```

Expected output:
```
🔍 Testing MeiliSearch initialization...
✅ MeiliSearch initialization test completed!
🎉 MeiliSearch is ready to use!
```

## Rollback

To rollback to Bleve search engine:

1. Set `MeiliSearch.Enabled: false` in `bbs-go.yaml`
2. Restart the application
3. The system will automatically use Bleve for search operations

## Performance Benefits

MeiliSearch offers several advantages over Bleve:

- **Faster search response times**
- **Better memory efficiency**
- **Advanced search features** (typo tolerance, synonyms, etc.)
- **Real-time updates**
- **Horizontal scalability**
- **RESTful API** for direct access if needed

## Troubleshooting

### Common Issues

1. **Connection Failed**: Ensure MeiliSearch is running on the configured host/port
2. **Authentication Error**: Verify the master key matches your MeiliSearch configuration
3. **Index Not Found**: The application will automatically create the index on first use
4. **Search Returns Empty**: Call `/api/search/reindex` to populate the search index

### Logs

Check application logs for MeiliSearch-related messages:
- `MeiliSearch initialized successfully` - Successful connection
- `MeiliSearch client not initialized` - Connection issues
- `Failed to add document to MeiliSearch` - Indexing errors

## Support

For issues or questions:
1. Check MeiliSearch server logs
2. Verify configuration settings
3. Test connection with `curl http://localhost:7700/health`
4. Review application logs for error messages