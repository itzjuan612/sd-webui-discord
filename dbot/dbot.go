/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 11:06:01
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-22 00:21:54
 * @Description: file content
 */
package dbot

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	AppCommands      []*discordgo.ApplicationCommand
	RegisteredCommands    []*discordgo.ApplicationCommand
	Session         *discordgo.Session
	ServerID        string
	SlashHandlerMap map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func NewDiscordBot(token string, serverID string) (*DiscordBot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	session.Client.Timeout = 120 * time.Second
	dbot := &DiscordBot{
		Session:         session,
		ServerID:        serverID,
		SlashHandlerMap: make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)),
		AppCommands:      make([]*discordgo.ApplicationCommand, 0),
	}

	// 预存长选项
	dbot.SetLongChoice()

	// 生成命令列表
	dbot.GenerateCommandList()
	dbot.SetLocation()

	// 根据命令列表生成Handler
	dbot.GenerateSlashMap()

	dbot.Session.AddHandler(dbot.Ready)
	dbot.Session.AddHandler(dbot.InteractionCreate)

	return dbot, nil
}

func (d *DiscordBot) Run() {
	err := d.Session.Open()
	if err != nil {
		log.Println(err)
		return
	}
	d.SyncCommands()
	defer d.Session.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
	//d.RemoveCommand()
	log.Println("Gracefully shutting down.")
}
