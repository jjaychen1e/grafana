import { EventBusSrv, BusEventWithPayload, EventBusGroup } from './EventBus';

interface LoginEventPayload {
  logins: number;
}

interface HelloEventPayload {
  hellos: number;
}
class LoginEvent extends BusEventWithPayload<LoginEventPayload> {
  public static type = 'login-event';
  public type = LoginEvent.type;
}

class HelloEvent extends BusEventWithPayload<HelloEventPayload> {
  public static type = 'hello-event';
  public type = HelloEvent.type;
}

describe('EventBus', () => {
  it('Can subscribe specific event', () => {
    const bus = new EventBusSrv();
    const events: LoginEvent[] = [];

    bus.on(LoginEvent, event => {
      events.push(event);
    });

    bus.emit(new LoginEvent({ logins: 10 }));
    bus.emit(new HelloEvent({ hellos: 10 }));

    expect(events[0].payload.logins).toBe(10);
    expect(events.length).toBe(1);
  });

  // it('Can subscribe to all events', () => {
  //   const bus = new EventBusSrv();
  //   let gotLoginEvent: LoginEvent | null = null;
  //   let gotHelloEvent: HelloEvent | null = null;
  //
  //   bus.subscribe((event: LoginEvent | HelloEvent) => {
  //     switch (event.type) {
  //       case LoginEvent.type: {
  //         gotLoginEvent = event;
  //       }
  //       case HelloEvent.type: {
  //         gotHelloEvent = event;
  //       }
  //     }
  //   });
  //
  //   bus.emit(new LoginEvent({ logins: 10 }));
  //   bus.emit(new HelloEvent({ hellos: 20 }));
  //
  //   expect(gotLoginEvent.payload.logins).toBe(10);
  //   expect(gotHelloEvent.payload.hellos).toBe(20);
  // });

  it('New group handles unsub', () => {
    const bus = new EventBusSrv();
    const group = new EventBusGroup(bus);
    const events: LoginEvent[] = [];

    group.on(LoginEvent, event => {
      events.push(event);
    });

    bus.emit(new LoginEvent({ logins: 10 }));

    expect(events.length).toBe(1);

    group.unsubscribe();

    bus.emit(new LoginEvent({ logins: 10 }));
    expect(events.length).toBe(1);
  });
});
