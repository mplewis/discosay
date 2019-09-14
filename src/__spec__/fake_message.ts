import * as Discord from 'discord.js'

export default function fakeMessage (content: string): Discord.Message {
  return { content } as Discord.Message
}
