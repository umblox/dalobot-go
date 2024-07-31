#!/bin/sh /etc/rc.common
# Copyright (C) 2024 arneta.id

START=99

start() {
    echo "Starting BOT"
    /usr/bin/telebot > /var/log/telebot.log 2>&1 &
}

stop() {
    echo "Stopping BOT"
    killall telebot
}

restart() {
    stop
    sleep 1
    start
}
