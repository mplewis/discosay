import { Responder } from '../types'
import crab from '../crab'

const trigger = 'riir'

const riir: Responder = {
  name: 'riir',
  applicable: msg => msg.content.toLowerCase() === trigger,
  handle: _ => crab('REWRITE IT IN RUST')
}

export default riir
