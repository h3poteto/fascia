{% extends "../layouts/webview.html.tpl" %}

{% block content %}
<div class="session">
  <div class="main">
    <div class="sign-in-board">
      <a href={{ publicURL }}><span class="pure-button button-success session-button"><span class="octicon octicon-mark-github"></span> Sign In with Github</span></a>
      <span class="message">This service does not access your private repository. If you want to management private repository, please click the button below.</span>
      <a href={{ privateURL }}><span class="pure-button button-small pure-button-primary secondary-session-button"><span class="octicon octicon-mark-github"></span> Sign In with Github Private Access</span></a>
    </div>
  </div>
</div>
{% endblock %}
