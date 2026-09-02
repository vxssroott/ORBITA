# ORBITA System Architecture

ORBITA is the independent software intelligence and operational layer between national space infrastructure and operational consumers.

                         SPACECRAFT
                             |
                             v
                    VENDOR / RF SYSTEMS
                             |
                             v
                 +-----------------------+
                 |        ORBITA         |
                 |                       |
                 | Telemetry Ingestion   |
                 | State Engine          |
                 | Event Engine          |
                 | Rules Engine          |
                 | Mission Operations    |
                 | Command Gateway       |
                 | Digital Twin          |
                 | Intelligence          |
                 | Security / Audit      |
                 +-----------------------+
                             |
             +---------------+---------------+
             |               |               |
             v               v               v
         OPERATORS        SYSTEMS        ANALYTICS

ORBITA does not need to replace spacecraft manufacturers' systems.

It provides an independent software intelligence and operational layer above them.
