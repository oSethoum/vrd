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
# vrd config file

# input vuerd file path
input: vrd/input.vuerd.json

# debug mode prints information
debug: true

# support sqlite3, mysql, postgres
database: sqlite3

# ent config
ent:
  output: ./output
  package: app
  privacy: true
  auth: true
  graphql:
    file_upload: true
    subscription: true
    relay_connection: true
```

## Docs

```yaml
# global options separator: |

Table:
   Name:
      nameMixin                  # decalaration of Mixin
   Comment:
      nm=users                   # entsql.Annotation{Table:"users"}
      mxs=(time, ..)             # TimeMixin{} in schema mixins
      -m2m                        # declare that table is m2m relationship
      nx=time                    # Mixin Alias to be referenced by tables in mxs(nx,...)
      skip=(mutation,subscription)
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
      skip=(op1,op2...):         # entgql.Skip(opt1,opt2...)
         all,                    # entgql.SkipAll()
         type,                   # entgql.SkipType()
         where                   # entgql.SkipWhereInput()
         create                  # entgql.SkipMutationCreateInput()
         update,                 # entgql.SkipMutationUpdateInput()
      upd=value                  # UpdateDefault(value)
      minLen=value               # MinLen(value)
      maxLen=value               # MaxLen(value)
      min=value                  # Min(value)
      max=value                  # Max(value)
      match=(regx)               # Match(regx)
      range=(x,y)                # Range(x,y)
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
         nr=(owner,todos)           # edge.To("todos", Todo.Type) | edge.form("owner", User.Type).Ref("todos")
         cr=name                    # indicate name in m2m middle table
         # case of m2n relationship with extra fields
         nr=(From,users,User,todos,through)   # edge.From("users", User.Type).Ref("todos").Through("tableName", TableName.Type)
         nr=(To,todos,Todo,through)           # edge.To("todos", Todo.Type).Through("tableName", TableName.Type)
```

## Supported ORMs

We currently support only [Ent](https://entgo.io/) but we are planning to support more (gorm, prisma) int the future.

## Contributions

we don't accept contributions at the moment, this is meant for personal use, but you are welcome to use it.
