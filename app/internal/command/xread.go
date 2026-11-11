package command

// func (h *Handler) handleXReadCommand(writer *bufio.Writer, args []string) error {
// 	if len(args) < 4 {
// 		return errors.New("invalid command")
// 	}

// 	output, err := h.store.ReadStream(args[1], args[2], args[3])
// 	if err != nil {
// 		return err
// 	}
// 	outputStr, err := SerializeStreamEntries(output)
// 	if err != nil {
// 		return err
// 	}
// 	writer.Write([]byte(outputStr))
// 	writer.Flush()
// 	return nil
// }

// func SerializeStreamEntry(entry storage.StreamEntry) (string, error) {

// }
