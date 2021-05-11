package db

type PooledDataSourceFactory struct {
	dataSources map[string]DataSource
}

func NewFactory() *PooledDataSourceFactory {
	var dataSourceFactory = new(PooledDataSourceFactory)
	dataSourceFactory.dataSources = make(map[string]DataSource)

	return dataSourceFactory
}

func (dataSourceFactory *PooledDataSourceFactory) GetDataSource(connectString string) (dataSource *DataSource, err error) {
	var td = dataSourceFactory.dataSources[connectString]
	dataSource = &td
	return
}
