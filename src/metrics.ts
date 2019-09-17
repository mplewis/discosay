import express from 'express'
import prom from 'prom-client'

const { PORT } = process.env
const PATH = '/metrics'

prom.collectDefaultMetrics()

const responses = new prom.Counter({
  name: 'responses',
  help: 'Messages the bot responded to',
  labelNames: ['username', 'responder']
})

export function responseSent (username: string, responder: string): void {
  responses.inc()
  responses.inc({ username })
  responses.inc({ responder })
}

export function serve (): void {
  const server = express()

  server.get(PATH, (_, res) => {
    res.set('Content-Type', prom.register.contentType)
    res.end(prom.register.metrics())
  })

  server.listen(PORT)
  console.log(`Serving metrics via ${PATH} on port ${PORT}`)
}
