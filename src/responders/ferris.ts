import { Responder } from '../types'
import crab from '../crab'

const trigger = '!retf '

const ferris: Responder = {
  name: 'ferris',
  applicable: msg => msg.content.toLowerCase().startsWith(trigger),
  handle: msg => crab(msg.content.slice(trigger.length))
}

export default ferris
