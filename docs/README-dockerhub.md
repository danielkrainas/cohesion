# Cohesion

Cohesion adds a lightweight, rudimentry service registry abstraction to [Serf](https://github.com/hashicorp/serf) agents and enables autonomically bootstrapping new nodes without prior knowledge of the cluster or in the event of a cluster failure. The service registry is managed through an HTTP API and saves registration information within labels in Serf. 

For more information, see the [project site.](https://github.com/danielkrainas/cohesion)

## Usage

```
$ docker run dakr/cohesion agent
```
