import { Responder } from '../types'

const ponger: Responder = {
  name: 'ponger',
  applicable: msg => msg.content.toLowerCase().includes('ping'),
  handle (msg) {
    return `${msg.content}: pong!`
  }
}

export default ponger
