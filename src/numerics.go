/* vim: set sw=4 sts=4 et foldmethod=syntax : */

/*
 * Copyright (c) 2011 Alexander Færøy <ahf@0x90.dk>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * * Redistributions of source code must retain the above copyright notice, this
 *   list of conditions and the following disclaimer.
 *
 * * Redistributions in binary form must reproduce the above copyright notice,
 *   this list of conditions and the following disclaimer in the documentation
 *   and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package main

import (
    "fmt"
)

type Numeric int

func Format(numeric Numeric, a ...interface{}) string {
    return fmt.Sprintf(replies[numeric], a...)
}

const (
    RPL_WELCOME Numeric     = 1
    RPL_YOURHOST            = 2
    RPL_CREATED             = 3
    RPL_MYINFO              = 4
    RPL_ISUPPORT            = 5
    RPL_REDIR               = 10
    RPL_MAP                 = 15 /* Undernet Extension. */
    RPL_MAPMORE             = 16 /* Undernet Extension. */
    RPL_MAPEND              = 17 /* Undernet Extension. */
    RPL_SAVENICK            = 43 /* From IRCnet.        */

    /* Numeric replies from server commands. */
    /* Currently in the range 200 to 399.    */
    RPL_TRACELINK           = 200
    RPL_TRACECONNECTING     = 201
    RPL_TRACEHANDSHAKE      = 202
    RPL_TRACEUNKNOWN        = 203
    RPL_TRACEOPERATOR       = 204
    RPL_TRACEUSER           = 205
    RPL_TRACESERVER         = 206
    RPL_TRACENEWTYPE        = 208
    RPL_TRACECLASS          = 209
    RPL_STATSLINKINFO       = 211
    RPL_STATSCOMMANDS       = 212
    RPL_STATSCLINE          = 213
    RPL_STATSNLINE          = 214
    RPL_STATSILINE          = 215
    RPL_STATSKLINE          = 216
    RPL_STATSQLINE          = 217
    RPL_STATSYLINE          = 218
    RPL_ENDOFSTATS          = 219
    RPL_STATSPLINE          = 220
    RPL_UMODEIS             = 221
    RPL_STATSFLINE          = 224
    RPL_STATSDLINE          = 225
    RPL_SERVLIST            = 234
    RPL_SERVLISTEND         = 235
    RPL_STATSLLINE          = 241
    RPL_STATSUPTIME         = 242
    RPL_STATSOLINE          = 243
    RPL_STATSHLINE          = 244
    RPL_STATSSLINE          = 245
    RPL_STATSXLINE          = 247
    RPL_STATSULINE          = 248
    RPL_STATSDEBUG          = 249
    RPL_STATSCONN           = 250
    RPL_LUSERCLIENT         = 251
    RPL_LUSEROP             = 252
    RPL_LUSERUNKNOWN        = 253
    RPL_LUSERCHANNELS       = 254
    RPL_LUSERME             = 255
    RPL_ADMINME             = 256
    RPL_ADMINLOC1           = 257
    RPL_ADMINLOC2           = 258
    RPL_ADMINEMAIL          = 259
    RPL_TRACELOG            = 261
    RPL_ENDOFTRACE          = 262
    RPL_LOAD2HI             = 263
    RPL_LOCALUSERS          = 265
    RPL_GLOBALUSERS         = 266
    RPL_ACCEPTLIST          = 281
    RPL_ENDOFACCEPT         = 282
    RPL_NONE                = 300
    RPL_AWAY                = 301
    RPL_USERHOST            = 302
    RPL_ISON                = 303
    RPL_TEXT                = 304
    RPL_UNAWAY              = 305
    RPL_NOWAWAY             = 306
    RPL_WHOISUSER           = 311
    RPL_WHOISSERVER         = 312
    RPL_WHOISOPERATOR       = 313
    RPL_WHOWASUSER          = 314
    RPL_ENDOFWHO            = 315
    RPL_ENDOFWHOWAS         = 369
    RPL_WHOISCHANOP         = 316
    RPL_WHOISIDLE           = 317
    RPL_ENDOFWHOIS          = 318
    RPL_WHOISCHANNELS       = 319
    RPL_LISTSTART           = 321
    RPL_LIST                = 322
    RPL_LISTEND             = 323
    RPL_CHANNELMODEIS       = 324
    RPL_CREATIONTIME        = 329
    RPL_WHOISLOGGEDIN       = 330
    RPL_NOTOPIC             = 331
    RPL_TOPIC               = 332
    RPL_TOPICWHOTIME        = 333
    RPL_WHOISACTUALLY       = 338
    RPL_INVITING            = 341
    RPL_SUMMONING           = 342
    RPL_INVITELIST          = 346
    RPL_ENDOFINVITELIST     = 347
    RPL_EXCEPTLIST          = 348
    RPL_ENDOFEXCEPTLIST     = 349
    RPL_VERSION             = 351
    RPL_WHOREPLY            = 352
    RPL_NAMREPLY            = 353
    RPL_KILLDONE            = 361
    RPL_CLOSING             = 362
    RPL_CLOSEEND            = 363
    RPL_LINKS               = 364
    RPL_ENDOFLINKS          = 365
    RPL_ENDOFNAMES          = 366
    RPL_BANLIST             = 367
    RPL_ENDOFBANLIST        = 368
    RPL_INFO                = 371
    RPL_MOTD                = 372
    RPL_INFOSTART           = 373
    RPL_ENDOFINFO           = 374
    RPL_MOTDSTART           = 375
    RPL_ENDOFMOTD           = 376
    RPL_YOUREOPER           = 381
    RPL_REHASHING           = 382
    RPL_MYPORTIS            = 384
    RPL_NOTOPERANYMORE      = 385
    RPL_RSACHALLENGE        = 386
    RPL_TIME                = 391
    RPL_USERSSTART          = 392
    RPL_USERS               = 393
    RPL_ENDOFUSERS          = 394
    RPL_NOUSERS             = 395

    /* Errors are in the range from 400 to 599. */
    ERR_NOSUCHNICK          = 401
    ERR_NOSUCHSERVER        = 402
    ERR_NOSUCHCHANNEL       = 403
    ERR_CANNOTSENDTOCHAN    = 404
    ERR_TOOMANYCHANNELS     = 405
    ERR_WASNOSUCHNICK       = 406
    ERR_TOOMANYTARGETS      = 407
    ERR_NOORIGIN            = 409
    ERR_INVALIDCAPCMD       = 410
    ERR_NORECIPIENT         = 411
    ERR_NOTEXTTOSEND        = 412
    ERR_NOTOPLEVEL          = 413
    ERR_WILDTOPLEVEL        = 414
    ERR_TOOMANYMATCHES      = 416
    ERR_UNKNOWNCOMMAND      = 421
    ERR_NOMOTD              = 422
    ERR_NOADMININFO         = 423
    ERR_FILEERROR           = 424
    ERR_NONICKNAMEGIVEN     = 431
    ERR_ERRONEUSNICKNAME    = 432
    ERR_NICKNAMEINUSE       = 433
    ERR_NICKCOLLISION       = 436
    ERR_UNAVAILRESOURCE     = 437
    ERR_NICKTOOFAST         = 438
    ERR_USERNOTINCHANNEL    = 441
    ERR_NOTONCHANNEL        = 442
    ERR_USERONCHANNEL       = 443
    ERR_NOLOGIN             = 444
    ERR_SUMMONDISABLED      = 445
    ERR_USERSDISABLED       = 446
    ERR_NOTREGISTERED       = 451
    ERR_ACCEPTFULL          = 456
    ERR_ACCEPTEXIST         = 457
    ERR_ACCEPTNOT           = 458
    ERR_NEEDMOREPARAMS      = 461
    ERR_ALREADYREGISTRED    = 462
    ERR_NOPERMFORHOST       = 463
    ERR_PASSWDMISMATCH      = 464
    ERR_YOUREBANNEDCREEP    = 465
    ERR_YOUWILLBEBANNED     = 466
    ERR_KEYSET              = 467
    ERR_CHANNELISFULL       = 471
    ERR_UNKNOWNMODE         = 472
    ERR_INVITEONLYCHAN      = 473
    ERR_BANNEDFROMCHAN      = 474
    ERR_BADCHANNELKEY       = 475
    ERR_BADCHANMASK         = 476
    ERR_NEEDREGGEDNICK      = 477
    ERR_BANLISTFULL         = 478
    ERR_BADCHANNAME         = 479
    ERR_SSLONLYCHAN         = 480
    ERR_NOPRIVILEGES        = 481
    ERR_CHANOPRIVSNEEDED    = 482
    ERR_CANTKILLSERVER      = 483
    ERR_ISCHANSERVICE       = 484
    ERR_BANNEDNICK          = 485
    ERR_VOICENEEDED         = 489
    ERR_NOOPERHOST          = 491
    ERR_UMODEUNKNOWNFLAG    = 501
    ERR_USERSDONTMATCH      = 502
    ERR_GHOSTEDCLIENT       = 503
    ERR_USERNOTONSERV       = 504
    ERR_WRONGPONG           = 513
    ERR_HELPNOTFOUND        = 524

    RPL_WHOISSECURE         = 671
    RPL_MODLIST             = 702
    RPL_ENDOFMODLIST        = 703
    RPL_HELPSTART           = 704
    RPL_HELPTXT             = 705
    RPL_ENDOFHELP           = 706
    ERR_TARGCHANGE          = 707
    RPL_ETRACEFULL          = 708
    RPL_ETRACE              = 709
    RPL_KNOCK               = 710
    RPL_KNOCKDLVR           = 711
    ERR_TOOMANYKNOCK        = 712
    ERR_CHANOPEN            = 713
    ERR_KNOCKONCHAN         = 714
    ERR_KNOCKDISABLED       = 715
    ERR_TARGUMODEG          = 716
    RPL_TARGNOTIFY          = 717
    RPL_UMODEGMSG           = 718
    RPL_OMOTDSTART          = 720
    RPL_OMOTD               = 721
    RPL_ENDOFOMOTD          = 722
    ERR_NOPRIVS             = 723
    RPL_TESTMASK            = 724
    RPL_TESTLINE            = 725
    RPL_NOTESTLINE          = 726
    RPL_TESTMASKGECOS       = 727
    RPL_MONONLINE           = 730
    RPL_MONOFFLINE          = 731
    RPL_MONLIST             = 732
    RPL_ENDOFMONLIST        = 733
    ERR_MONLISTFULL         = 734
    RPL_RSACHALLENGE2       = 740
    RPL_ENDOFRSACHALLENGE2  = 741
)

var replies = map[Numeric] string {
    RPL_WELCOME:            ":%s 001 %s :Welcome to the %s Internet Relay Chat Network %s",
    RPL_YOURHOST:           ":%s 002 %s :Your host is %s, running version %s",
    RPL_CREATED:            ":%s 003 %s :This server was created %s",
    RPL_MYINFO:             ":%s 004 %s %s %s oiwszcrkfydnxbauglZCD biklmnopstveIrS bkloveI",
    RPL_ISUPPORT:           "%s :are supported by this server",
    RPL_REDIR:              ":%s 010 %s %s %d :Please use this Server/Port instead",
    RPL_MAP:                ":%s 015 %s :%s",
    RPL_MAPEND:             ":%s 017 %s :End of /MAP",
    RPL_SAVENICK:           "%s :Nick collision, forcing nick change to your unique ID",
    RPL_TRACELINK:          "Link %s %s %s",
    RPL_TRACECONNECTING:    "Try. %s %s",
    RPL_TRACEHANDSHAKE:     "H.S. %s %s",
    RPL_TRACEUNKNOWN:       "???? %s %s (%s) %d",
    RPL_TRACEOPERATOR:      "Oper %s %s (%s) %lu %lu",
    RPL_TRACEUSER:          "User %s %s (%s) %lu %lu",
    RPL_TRACESERVER:        "Serv %s %dS %dC %s %s!%s@%s %lu",
    RPL_TRACENEWTYPE:       "<newtype> 0 %s",
    RPL_TRACECLASS:         "Class %s %d",
    RPL_STATSCOMMANDS:      "%s %u %u :%u",
    RPL_STATSCLINE:         "C %s %s %s %d %s",
    RPL_STATSILINE:         "I %s * %s@%s %d %s",
    RPL_STATSKLINE:         "%c %s * %s :%s%s%s",
    RPL_STATSQLINE:         "%c %d %s :%s",
    RPL_STATSYLINE:         "Y %s %d %d %d %u %d.%d %d.%d %u",
    RPL_ENDOFSTATS:         "%c :End of /STATS report",
    RPL_STATSPLINE:         "%c %d %s %d :%s%s",
    RPL_UMODEIS:            ":%s 221 %s %s",
    RPL_STATSDLINE:         "%c %s :%s%s%s",
    RPL_STATSLLINE:         "L %s * %s 0 -1",
    RPL_STATSUPTIME:        ":Server Up %d days, %d:%02d:%02d",
    RPL_STATSOLINE:         "O %s@%s * %s %s %s",
    RPL_STATSHLINE:         "H %s * %s 0 -1",
    RPL_STATSXLINE:         "%c %d %s :%s",
    RPL_STATSULINE:         "U %s %s@%s %s",
    RPL_STATSCONN:          ":Highest connection count: %d (%d clients) (%d connections received)",
    RPL_LUSERCLIENT:        ":There are %d users and %d invisible on %d servers",
    RPL_LUSEROP:            "%d :IRC Operators online",
    RPL_LUSERUNKNOWN:       "%d :unknown connection(s)",
    RPL_LUSERCHANNELS:      "%d :channels formed",
    RPL_LUSERME:            ":I have %d clients and %d servers",
    RPL_ADMINME:            ":%s 256 %s :Administrative info about %s",
    RPL_ADMINLOC1:          ":%s 257 %s :%s",
    RPL_ADMINLOC2:          ":%s 258 %s :%s",
    RPL_ADMINEMAIL:         ":%s 259 %s :%s",
    RPL_ENDOFTRACE:         "%s :End of TRACE",
    RPL_LOAD2HI:            ":%s 263 %s %s :Server load is temporarily too heavy. Please wait a while and try again.",
    RPL_LOCALUSERS:         "%d %d :Current local users %d, max %d",
    RPL_GLOBALUSERS:        "%d %d :Current global users %d, max %d",
    RPL_ACCEPTLIST:         ":%s 281 %s %s",
    RPL_ENDOFACCEPT:        ":%s 282 %s :End of /ACCEPT list.",
    RPL_AWAY:               "%s :%s",
    RPL_USERHOST:           ":%s 302 %s :%s",
    RPL_ISON:               ":%s 303 %s :",
    RPL_UNAWAY:             ":%s 305 %s :You are no longer marked as being away",
    RPL_NOWAWAY:            ":%s 306 %s :You have been marked as being away",
    RPL_WHOISUSER:          "%s %s %s * :%s",
    RPL_WHOISSERVER:        "%s %s :%s",
    RPL_WHOISOPERATOR:      "%s :%s",
    RPL_WHOWASUSER:         ":%s 314 %s %s %s %s * :%s",
    RPL_ENDOFWHO:           ":%s 315 %s %s :End of /WHO list.",
    RPL_WHOISIDLE:          "%s %d %d :seconds idle, signon time",
    RPL_ENDOFWHOIS:         "%s :End of /WHOIS list.",
    RPL_WHOISCHANNELS:      ":%s 319 %s %s :",
    RPL_LISTSTART:          ":%s 321 %s Channel :Users  Name",
    RPL_LIST:               ":%s 322 %s %s %d :%s",
    RPL_LISTEND:            ":%s 323 %s :End of /LIST",
    RPL_CHANNELMODEIS:      ":%s 324 %s %s %s",
    RPL_CREATIONTIME:       ":%s 329 %s %s %lu",
    RPL_WHOISLOGGEDIN:      ":%s 330 %s %s %s :is logged in as",
    RPL_NOTOPIC:            ":%s 331 %s %s :No topic is set.",
    RPL_TOPIC:              ":%s 332 %s %s :%s",
    RPL_TOPICWHOTIME:       ":%s 333 %s %s %s %lu",
    RPL_WHOISACTUALLY:      "%s %s :actually using host",
    RPL_INVITING:           ":%s 341 %s %s %s",
    RPL_EXCEPTLIST:         ":%s 348 %s %s %s %s %lu",
    RPL_ENDOFEXCEPTLIST:    ":%s 349 %s %s :End of Channel Exception List",
    RPL_VERSION:            "%s(%s). %s :%s TS%dow %s",
    RPL_WHOREPLY:           ":%s 352 %s %s %s %s %s %s %s :%d %s",
    RPL_NAMREPLY:           ":%s 353 %s %s %s :",
    RPL_CLOSING:            ":%s 362 %s %s :Closed. Status = %d",
    RPL_CLOSEEND:           ":%s 363 %s %d :Connections Closed",
    RPL_LINKS:              "%s %s :%d %s",
    RPL_ENDOFLINKS:         "%s :End of /LINKS list.",
    RPL_ENDOFNAMES:         ":%s 366 %s %s :End of /NAMES list.",
    RPL_BANLIST:            ":%s 367 %s %s %s %s %lu",
    RPL_ENDOFBANLIST:       ":%s 368 %s %s :End of Channel Ban List",
    RPL_ENDOFWHOWAS:        ":%s 369 %s %s :End of WHOWAS",
    RPL_INFO:               ":%s",
    RPL_MOTD:               ":%s 372 %s :- %s",
    RPL_ENDOFINFO:          ":End of /INFO list.",
    RPL_MOTDSTART:          ":%s 375 %s :- %s Message of the Day - ",
    RPL_ENDOFMOTD:          ":%s 376 %s :End of /MOTD command.",
    RPL_YOUREOPER:          ":%s 381 %s :You have entered... the PANTS FREE ZONE!",
    RPL_REHASHING:          ":%s 382 %s %s :Rehashing",
    RPL_RSACHALLENGE:       ":%s 386 %s :%s",
    RPL_TIME:               "%s :%s",
    ERR_NOSUCHNICK:         "%s :No such nick/channel",
    ERR_NOSUCHSERVER:       "%s :No such server",
    ERR_NOSUCHCHANNEL:      "%s :No such channel",
    ERR_CANNOTSENDTOCHAN:   "%s :Cannot send to channel",
    ERR_TOOMANYCHANNELS:    ":%s 405 %s %s :You have joined too many channels",
    ERR_WASNOSUCHNICK:      ":%s 406 %s %s :There was no such nickname",
    ERR_TOOMANYTARGETS:     ":%s 407 %s %s :Too many recipients.",
    ERR_NOORIGIN:           ":%s 409 %s :No origin specified",
    ERR_INVALIDCAPCMD:      ":%s 410 %s %s :Invalid CAP subcommand",
    ERR_NORECIPIENT:        ":%s 411 %s :No recipient given (%s)",
    ERR_NOTEXTTOSEND:       ":%s 412 %s :No text to send",
    ERR_NOTOPLEVEL:         "%s :No toplevel domain specified",
    ERR_WILDTOPLEVEL:       "%s :Wildcard in toplevel Domain",
    ERR_TOOMANYMATCHES:     ":%s 416 %s %s :output too large, truncated",
    ERR_UNKNOWNCOMMAND:     ":%s 421 %s %s :Unknown command",
    ERR_NOMOTD:             ":%s 422 %s :MOTD File is missing",
    ERR_NONICKNAMEGIVEN:    ":%s 431 %s :No nickname given",
    ERR_ERRONEUSNICKNAME:   ":%s 432 %s %s :Erroneous Nickname",
    ERR_NICKNAMEINUSE:      ":%s 433 %s %s :Nickname is already in use.",
    ERR_NICKCOLLISION:      "%s :Nickname collision KILL",
    ERR_UNAVAILRESOURCE:    ":%s 437 %s %s :Nick/channel is temporarily unavailable",
    ERR_NICKTOOFAST:        ":%s 438 %s %s %s :Nick change too fast. Please wait %d seconds.",
    ERR_USERNOTINCHANNEL:   "%s %s :They aren't on that channel",
    ERR_NOTONCHANNEL:       "%s :You're not on that channel",
    ERR_USERONCHANNEL:      "%s %s :is already on channel",
    ERR_NOTREGISTERED:      ":%s 451 * :You have not registered",
    ERR_ACCEPTFULL:         ":%s 456 %s :Accept list is full",
    ERR_ACCEPTEXIST:        ":%s 457 %s %s :is already on your accept list",
    ERR_ACCEPTNOT:          ":%s 458 %s %s :is not on your accept list",
    ERR_NEEDMOREPARAMS:     ":%s 461 %s %s :Not enough parameters",
    ERR_ALREADYREGISTRED:   ":%s 462 %s :You may not reregister",
    ERR_PASSWDMISMATCH:     ":%s 464 %s :Password Incorrect",
    ERR_YOUREBANNEDCREEP:   ":%s 465 %s :You are banned from this server- %s",
    ERR_CHANNELISFULL:      ":%s 471 %s %s :Cannot join channel (+l)",
    ERR_UNKNOWNMODE:        ":%s 472 %s %c :is unknown mode char to me",
    ERR_INVITEONLYCHAN:     ":%s 473 %s %s :Cannot join channel (+i)",
    ERR_BANNEDFROMCHAN:     ":%s 474 %s %s :Cannot join channel (+b)",
    ERR_BADCHANNELKEY:      ":%s 475 %s %s :Cannot join channel (+k)",
    ERR_NEEDREGGEDNICK:     ":%s 477 %s %s :Cannot join channel (+r)",
    ERR_BANLISTFULL:        ":%s 478 %s %s %s :Channel ban list is full",
    ERR_BADCHANNAME:        "%s :Illegal channel name",
    ERR_SSLONLYCHAN:        ":%s 480 %s %s :Cannot join channel (+S) - SSL/TLS required",
    ERR_NOPRIVILEGES:       ":Permission Denied - You're not an IRC operator",
    ERR_CHANOPRIVSNEEDED:   ":%s 482 %s %s :You're not channel operator",
    ERR_CANTKILLSERVER:     ":You can't kill a server!",
    ERR_ISCHANSERVICE:      ":%s 484 %s %s %s :Cannot kick or deop a network service",
    ERR_VOICENEEDED:        ":%s 489 %s %s :You're neither voiced nor channel operator",
    ERR_NOOPERHOST:         ":%s 491 %s :Only few of mere mortals may try to enter the twilight zone",
    ERR_UMODEUNKNOWNFLAG:   ":%s 501 %s :Unknown MODE flag",
    ERR_USERSDONTMATCH:     ":%s 502 %s :Can't change mode for other users",
    ERR_USERNOTONSERV:      ":%s 504 %s %s :User is not on this server",
    ERR_WRONGPONG:          ":%s 513 %s :To connect type /QUOTE PONG %08lX",
    ERR_HELPNOTFOUND:       ":%s 524 %s %s :Help not found",
    RPL_WHOISSECURE:        "%s :is using a secure connection",
    RPL_MODLIST:            ":%s 702 %s %s 0x%p %s %s",
    RPL_ENDOFMODLIST:       ":%s 703 %s :End of /MODLIST.",
    RPL_HELPSTART:          ":%s 704 %s %s :%s",
    RPL_HELPTXT:            ":%s 705 %s %s :%s",
    RPL_ENDOFHELP:          ":%s 706 %s %s :End of /HELP.",
    ERR_TARGCHANGE:         ":%s 707 %s %s :Targets changing too fast, message dropped",
    RPL_ETRACEFULL:         ":%s 708 %s %s %s %s %s %s %s %s :%s",
    RPL_ETRACE:             ":%s 709 %s %s %s %s %s %s %s :%s",
    RPL_KNOCK:              ":%s 710 %s %s %s!%s@%s :has asked for an invite.",
    RPL_KNOCKDLVR:          ":%s 711 %s %s :Your KNOCK has been delivered.",
    ERR_TOOMANYKNOCK:       ":%s 712 %s %s :Too many KNOCKs (%s).",
    ERR_CHANOPEN:           "%s :Channel is open.",
    ERR_KNOCKONCHAN:        ":%s 714 %s %s :You are already on that channel.",
    ERR_KNOCKDISABLED:      ":%s 715 %s :KNOCKs are disabled.",
    ERR_TARGUMODEG:         "%s :is in +g mode (server-side ignore.)",
    RPL_TARGNOTIFY:         "%s :has been informed that you messaged them.",
    RPL_UMODEGMSG:          ":%s 718 %s %s %s@%s :is messaging you, and you have umode +g.",
    RPL_OMOTDSTART:         ":%s 720 %s :Start of OPER MOTD",
    RPL_OMOTD:              ":%s 721 %s :%s",
    RPL_ENDOFOMOTD:         ":%s 722 %s :End of OPER MOTD",
    ERR_NOPRIVS:            ":%s 723 %s %s :Insufficient oper privs",
    RPL_TESTLINE:           ":%s 725 %s %c %ld %s :%s",
    RPL_NOTESTLINE:         ":%s 726 %s %s :No matches",
    RPL_TESTMASKGECOS:      ":%s 727 %s %d %d %s!%s@%s %s :Local/remote clients match",
    RPL_MONONLINE:          ":%s 730 %s :%s",
    RPL_MONOFFLINE:         ":%s 731 %s :%s",
    RPL_MONLIST:            ":%s 732 %s :%s",
    RPL_ENDOFMONLIST:       ":%s 733 %s :End of MONITOR list",
    ERR_MONLISTFULL:        ":%s 734 %s %d %s :Monitor list is full",
    RPL_RSACHALLENGE2:      ":%s 740 %s :%s",
    RPL_ENDOFRSACHALLENGE2: ":%s 741 %s :End of CHALLENGE",
}
