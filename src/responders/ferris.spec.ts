import ferris from "./ferris";
import fakeMessage from '../../spec/fake_message'

describe('ferris', () => {
  it('applies to messages starting with the trigger', () => {
    expect(ferris.applicable(fakeMessage('!retf rust is great!'))).toBe(true);
    expect(ferris.applicable(fakeMessage('rust is great!'))).toBe(false);
    expect(ferris.applicable(fakeMessage('!retf'))).toBe(false);
  })
})
