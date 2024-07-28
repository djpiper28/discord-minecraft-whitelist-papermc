# minecraft-server
mc server bot and, plugin to allow for a compsoc whitelist

## User flow

### 1 Register User

```
/mcadd <minecraft username>
```

### 2 Verify User

The user joins the Miencraft Server, they will be kicked with a message showing a code.

Example Message:
```
You need to verify your account, use /mcverify 123456
```

### 3 User Verifies Themself

The user uses `/mcverify 123456`

### 4 The User Can Now Use The Server

YAY!!!!!!!!!

## Banning A User

On Minecraft, or Discord a user can be `/mcban`ed, the user is then kicked and not allowed to join
the Discord user and all of their alt accounts will be tagged as banned.

## Unbanning A User

Users can be unbanned by `/mcunban`ing someone on the Discord Bot, they will have all alt accounts freed up.
