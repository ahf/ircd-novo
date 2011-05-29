MISC_DIR=../misc

build_information.go:
	$(MISC_DIR)/buildinfo.sh --package main > $@

authors.go: ../AUTHORS.md
	$(MISC_DIR)/file_to_slice.sh --package main --name AuthorsFile --file $< > $@

include $(GOROOT)/src/Make.inc

TARG = ircd-novo
GOFILES = \
	authors.go \
	build_information.go \
	channel.go \
	channel_registry.go \
	channel_set.go \
	channel_utilities.go \
	client.go \
	client_registry.go \
	client_set.go \
	config.go \
	dns.go \
	hash.go \
	ircd.go \
	ircd_novo.go \
	join_handler.go \
	list_handler.go \
	listeners.go \
	match.go \
	message_handler.go \
	misc.go \
	motd_handler.go \
	notice_handler.go \
	notice_message.go \
	numerics.go \
	parser.go \
	part_handler.go \
	ping_handler.go \
	private_message.go \
	privmsg_handler.go \
	quit_handler.go \
	topic.go \
	topic_handler.go \
	topic_message.go \
	version.go \
	whois_handler.go \

CLEANFILES += authors.go build_information.go

include $(GOROOT)/src/Make.cmd
