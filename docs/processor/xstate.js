// Available variables:
// - Machine
// - interpret
// - assign
// - send
// - sendParent
// - spawn
// - raise
// - actions
// - XState (all XState exports)

const cacheStates = {
  initial: "undetermined",
  states: {
    undetermined: {
      on: {
        DATA: "determining"
      }
    },
    determining: {
      on: {
        SUCCESS: "determined",
        FAILURE: "undetermined"
      }
    },
    determined: {}
  }
};

const processorMachine = Machine({
  id: "processor",
  initial: "uninitialized",
  context: {
    state: 0,
    configuration: 0,
    data: []
  },
  states: {
    uninitialized: {
      on: {
        CONFIGURE_SUCCESS: "configured",
        CONFIGURE_FAILURE: "uninitialized"
      }
    },
    configured: {
      on: {
        DATA: "has_state"
      }
    },
    has_state: {
      on: {
        DONE: "running"
      },
      ...cacheStates
    },
    running: {
      on: {
        SUCCESS: "done",
        FAILURE: "configured"
      }
    },
    done: {
      type: "final"
    }
  }
});
