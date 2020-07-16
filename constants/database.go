package constants

// Database constants
const (
	// TO_DO: Replace this connection string with srv when this issue with go is resolved https://github.com/golang/go/issues/37362
	MongoConnectionString = "mongodb://lmsservice:tPC2VEbj7kcGucSr@cluster0-shard-00-00-sjgkt.gcp.mongodb.net:27017,cluster0-shard-00-01-sjgkt.gcp.mongodb.net:27017,cluster0-shard-00-02-sjgkt.gcp.mongodb.net:27017/test?ssl=true&replicaSet=Cluster0-shard-0&authSource=admin&retryWrites=true&w=majority"
	DatabaseName          = "lms_internal"
)

// Collection names
const (
	Leads         = "leads"
	LeadStatus    = "lead_status"
	Campaigns     = "campaigns"
	LeadTemplates = "lead_templates"
	Sources       = "sources"
)
