# Payment Platform Documentation

The canonical architecture and engineering reference is in [architecture.md](architecture.md).

It covers service ownership, financial invariants, database concurrency, idempotency, NATS JetStream, the transactional outbox, Saga workflows, FX/risk/reconciliation, security, observability, deployment and test operations.

The architecture includes six illustrations in [`assets/`](assets/). All image paths in `architecture.md` are relative to this directory and render on GitHub.

## Documentation maintenance

- Update the architecture document and affected figures together with changes to service boundaries, financial invariants or event contracts.
- Treat SQL migrations, API contracts and executable tests as authoritative for exact implementation behavior.
- Do not state that this reference backend is certified for production financial transactions without separate validation.
