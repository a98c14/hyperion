## Database Design
A GameEntity is composed of components. Each component has an id, name and possible set of values

```json
// Component
{
    "Id": 126,
    "Name": "BaseStats",
    "Parent": null,
    // "Children": [12, 42],
}
{
    "Id": 12,
    "Name": "OffensiveStats",
    "Parent": 126,
    // "Children": [2, 4, 6, 7],
}
{
    "Id": 2,
    "Name": "Power",
    "Parent": 12,
    // "Children": null,
}
{
    "Id": 3,
    "Name": "Haste",
    "Parent": 12,
    // "Children": null,
}

{
    "Id": 17,
    "Name": "Enemy",
    "Parent": null,
}

{
    "Id": 33,
    "Name": "SkillBuffer",
    "Parent": null,
}
```

Entities
```json
// Prefabs
{
    "Id": 12,
    "Name": "EnemyBoar",
    // "Components": [2, 3]
}

// PrefabComponents
{
    "PrefabId": 12,
    "ComponentId": 2,
    "Value": 50,
},
{
    "PrefabId": 12,
    "ComponentId": 3,
    "Value": 20,
}
```


