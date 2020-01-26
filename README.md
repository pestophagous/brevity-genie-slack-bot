> Tip: Render *.md files locally with Chrome browser http://stackoverflow.com/a/15626336

What purpose does this slack bot serve?
---------------------------------------

The following screenshot shows a typical slack conversation that you want
brevity-genie to interrupt for you:

![slack convo screenshot](artwork/sample_use.png)

Whenever you are in a channel that brevity-genie has joined, brevity-genie will
keep a count of how many messages you have typed in the past 30 minutes or so
(depending on configuration).

If you've been chatting for 7 minutes (by default), brevity-genie will jump in
and remind you that you must be getting back to work. At that point, you need to
allow 20 minutes (by default) to elapse with no more messages from you. After
such a silence period has elapsed, then the next time you type a message you
won't get any reminders from brevity-genie until you once again reach 7 minutes
of chatting.

The benefits of this approach:

 - genie tells everyone you must be going -- so you don't have to!
 - even if you're happy to sign off all by yourself, genie has better time-keeping.
 - some people are more willing to take commands from 'authorities' (which genie might be percieved as).
 - some people might be more polite to bots.

https://chatbotsmagazine.com/rethinking-politeness-in-chatbots-interesting-research-findings-6c3a91871e5f

Bare-bones, likely-deficient installation guide:
------------------------------------------------

This README is not the right place to fully document how to install and manage
slack bots in your slack workspace.

This is only a *rough* overview. You should read official Slack documents first.

However, since there are many attributes to each app you install, this
bare-bones guide can at least help you narrowly focus on which attributes to pay
attention to when installing brevity-genie.

### External Slack Documents ###

 - https://api.slack.com/start
 - https://slack.com/help/articles/115005265703-Create-a-bot-for-your-workspace
 - https://api.slack.com/tutorials/watson-sentiment
 - https://api.slack.com/tutorials/tunneling-with-ngrok

### Likely-deficient "installation steps" for brevity-genie ###

Go to https://api.slack.com/ then look for 'Your Apps'.

Click 'Create New App'. (Yes, the terminology is a bit weird. To my mind,
"creating the app" is what I was doing when I wrote all the Go code in this
repository. What actually happens when you click 'Create New App' is nothing
like coding the logic for an app. It really just creates KEYS and TOKENS and
SECRETS and PERMISSIONS to use for running your instance of the brevity-genie
codebase.)

After clicking 'Create New App', enter a name for this app instance, choose your
workspace as the development workspace, then click 'Create App'. Now you should
be in some sort of settings/configuration screen for your app.

Find your way to: Permissions >> Scopes >> Bot Token Scopes

Just add one for now: `app_mentions:read`

(We add more permissions later, but you need to add one in order to unlock
certain intermediate steps.)

Find your way to: Basic Information >> Display Information

In 'Display Information', add an app icon. **This will be the avatar for the bot
whenever the bot participates in conversations.**

Under "app credentials" retrieve the "Verification token". (It should look
approximately like `TBXoBKLwbs7WWWWWhbKM8QUq` and you will pass it to the code
as environment var SLACKVTOKEN_DEPRECATED)

You will need a bot token, too, but you cannot get that until you install the
app. (At least, such was the workflow at time of writing.)

Now install the app to your workspace.

Now go to permissions and you should find "Bot User OAuth Access Token". (It
should look approximately like
`xoxb-258044125036-017033000000-CB4Uo0Wg5cazN5UXXXXJB5bY` and you will pass it
to the code as environment var SLACKTOKEN)

Now on your server:
```
sudo setcap CAP_NET_BIND_SERVICE=+eip genie # <-- so you need not use root/sudo for the next line.
SLACKTOKEN=xoxb-258044125036-017033000000-CB4Uo0Wg5cazN5UXXXXJB5bY   SLACKVTOKEN_DEPRECATED=TBXoBKLwbs7WWWWWhbKM8QUq ./genie
```

Once your server is running:

Go back to the slack settings and look for 'Event Subscriptions' and the place
to 'Enable Events'. You will need to enter the URL where your server can be
reached. Slack needs to send a challenge message to verify the app. (The
brevity-genie code is designed to handle this challenge; all you need to do is
launch the app on your server as just shown above.)

If you are testing the app locally and do not have a publicly-reachable URL,
then read this tutorial: https://api.slack.com/tutorials/tunneling-with-ngrok

--------------------------------------------------------------------------------

![slack convo screenshot](artwork/event_subscriptions.png)

--------------------------------------------------------------------------------

Choose as many of these permissions as you can find on the Events page:

```
app_mentions:read
channels:history
channels:join
channels:read
chat:write
groups:history
im:history
im:read
im:write
mpim:history
mpim:read
mpim:write
users.profile:read
users:read
users:read.email
users:write
```

Save changes. Reinstall app. (Clicking on 'reinstall' is necessary every time
you add permissions.)

Now go to back to: Features >> Permissions >> Scopes

In this other permission area, you should be able to locate and add any of the
remaining permissions shown in the preceding list (i.e. any that you did not
find under the Events settings).



Acknowledgements
----------------

Thanks, [], for []

Thanks, [], which I use on all my Go code to make cool architectural diagrams like this one:
