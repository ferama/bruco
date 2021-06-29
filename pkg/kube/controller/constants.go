package kubecontroller

const controllerAgentName = "bruco-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Bruco is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Bruco fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Bruco"
	// MessageResourceSynced is the message used for an Event fired when a Bruco
	// is synced successfully
	MessageResourceSynced = "Bruco synced successfully"
)
