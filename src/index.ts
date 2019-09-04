import * as Discord from 'discord.js'
import ponger from './responders/ponger'

const responders = [ponger]

const { TOKEN } = process.env
if (!TOKEN) throw new Error('TOKEN is unset')

const client = new Discord.Client()

let name = '<unset>'
client.on('ready', () => {
  const { user } = client
  if (!user) throw new Error('No user on client')
  console.log(`Signed in as ${user.tag}`)
  name = user.username
})

client.on('message', msg => {
  const { author } = msg
  if (!author) return
  if (author.bot) return

  responders.forEach(responder => {
    if (!responder.applicable(msg)) return

    console.log(`${responder.name} ← ${author.username}: ${msg.content}`)
    const response = responder.handle(msg)
    console.log(`${responder.name} → ${name}: ${response}`)
    msg.channel.send(response)
  })
})

client.login(TOKEN)
