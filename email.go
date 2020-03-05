/*
	Email Message Body:
	1. from: source address e.g. "wilma@here.com"
	2. to: destination address string e.g. "betty@there.com"
	3. body: message content e.g. "Hi Betty, hope you're doing well."

	Email Server:
	1. Mail Submission Agent (MSA): moves user A's email into user A's outbox
	2. Mail Transfer Agent (MTA): uses MSA to read and delete message from user A's outbox, periodically.
	   These messages are sent to user B's email server
	3. User B's MTA uses its MSA to add the message to another user's inbox
	
	At any time, a user may ask the MSA to list the messages in their outbox
	At any time, a user may ask the MSA to list, read and delete the messages in their inbox
	The network address of an email server may be obtained by supplying the source or destination address
	of an email address to a Blue Book server
*/