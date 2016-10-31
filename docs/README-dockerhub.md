# Cohesion

Cohesion is a tool for automatically clustering nodes to one another. It will monitor the cluster agent on the node for when or if the node becomes disconnected or if it was never connected to start, Cohesion will use a preconfigured strategy for locating and joining the cluster or another stray node. The only supported node agent right now is [Hashicorp's Serf](https://github.com/hashicorp/serf).

For more information, see the [project site.](https://github.com/danielkrainas/cohesion)

## Usage

```
$ docker run dakr/cohesion agent
```
