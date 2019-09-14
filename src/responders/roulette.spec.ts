import roulette from "./roulette";
import fakeMessage from '../../spec/fake_message'

describe('roulette', () => {
  it('applies to messages matching the trigger', () => {
    expect(roulette.applicable(fakeMessage('!retf'))).toBe(true);
    expect(roulette.applicable(fakeMessage('!retf seriously, please use rust'))).toBe(false);
  })
})
