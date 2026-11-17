package command

// func (h *Handler) handleExecCommand(client *client.Client, args []string) {
// 	if len(args) > 1 {
// 		client.WriteResponse(errors.New("invalid command"))
// 		return
// 	}
// 	if !client.InTransaction {
// 		client.WriteResponse(errors.New("ERR EXEC without MULTI"))
// 		return
// 	}
// 	client.InTransaction = false
// 	responses := []any{}
// 	for args := range client.QueuedCommands {

// 	}
// 	return
// }
