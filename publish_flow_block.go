package rabbitmq

func (publisher *Publisher) startNotifyFlowHandler() {
	notifyFlowChan := publisher.chanManager.NotifyFlowSafe(make(chan bool))
	publisher.disablePublishDueToFlowMu.Lock()
	publisher.disablePublishDueToFlow = false
	publisher.disablePublishDueToFlowMu.Unlock()

	for ok := range notifyFlowChan {
		publisher.disablePublishDueToFlowMu.Lock()
		if ok {
			publisher.options.Logger.Warnf("pausing publishing due to flow request from server")
			publisher.disablePublishDueToFlow = true
		} else {
			publisher.disablePublishDueToFlow = false
			publisher.options.Logger.Warnf("resuming publishing due to flow request from server")
		}
		publisher.disablePublishDueToFlowMu.Unlock()
	}
}
