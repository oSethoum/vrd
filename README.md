# vrd

## What is vrd ?

The role of vuerd is to generate code for ORMs based on the [ERD Editor](https://marketplace.visualstudio.com/items?itemName=dineug.vuerd-vscode) vscode extension.

## Docs:

Download the binary and add it to your path, and from your command line

```console
vrd
```

that will initialize a config file `vrd.config.yaml`

```yaml
# vuerd schema path
input: vrd/db.vuerd.json

# files output path
output: ./output

# ent config
ent:
  package: app
  graphql: true
  auth: true
  privacy: true
  file_upload: true
  debug: true
  database: sqlite3 # sqlite3 | postgres | mysql
```

## Docs

```yaml
# global options separator: |

Table:
   Name:
      nameMixin                  # decalaration of Mixin
   Comment:
      nm=users                   # entsql.Annotation{Table:"users"}
      mxs(time, ..)              # TimeMixin{} in schema mixins
      m2m                        # declare that table is m2m relationship
      nx=time                    # Mixin Alias to be referenced by tables in mxs(nx,...)

Column:
   Name:                         # use: SnakeCase
   Unique:                       # Unique()
   AI:                           # AutoIncrement()
   NULL:                         # Optional().Nillable()
   DataType:
      Enum(toGo,done)            # Enum(Ae,Be).{Values(ToGo, Done)|NamedValues(ToGo,to_go, Done, done)}
   Default:
      time.Now                   # Default(time.Now)
   Comment:
      skip(op1,op2...):          # entgql.Skip(opt1,opt2...)
         options:
            create               # entgql.SkipMutationCreateInput()
            update,              # entgql.SkipMutationUpdateInput()
            type,                # entgql.SkipType()
            all,                 # entgql.SkipAll()
            where                # entgql.SkipWhereInput()
      upd=value                  # UpdateDefault(value)
      min=value                  # MinLength(value) or Min(value)
      max=value                  # MaxLength(value) or Max(value)
      match(regx)                # Match(regx)
      range(x,y)                 # Range(x,y)
      -nem                       # NotEmpty()
      -pos                       # Positive()
      -neg                       # Negative()
      -nneg                      # NonNegative()
      -im                        # Immutable()
      -s                         # Sensitive()
      -op                        # Optional()

Relationship:
   Type:
      OneN:
         start                   # edge.To()
         end                     # edge.From().Unqiue().Required()
      ZeroN:
         start                   # edge.To()
         end                     # edge.From().Unqiue()
      OneOnly:
         start                   # edge.To().Unique()
         end                     # edge.From().Unqiue().Required()
      ZeroOne:
         start                   # edge.To().Unique()
         end                     # edge.From().Unqiue()
   FK:
      Comment:
         # normal relationship
         (owner,todos)           # edge.To("todos", Todo.Type) | edge.form("owner", User.Type).Ref("todos")

         # case of m2n relationship with extra fields
         From,users,User,todos   # edge.From("users", User.Type).Ref("todos").Through("tableName", TableName.Type)
         To,todos,Todo           # edge.To("todos", Todo.Type).Through("tableName", TableName.Type)
```

## Supported ORMs

We currently support only [Ent](https://entgo.io/) but we are planning to support more (gorm, prisma) int the future.

## Contributions

we don't accept contributions at the moment, this is meant for personal use, but you are welcome to use it.
