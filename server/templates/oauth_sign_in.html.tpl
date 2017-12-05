{% extends "layouts/not_login.html.tpl" %}

{% block title %}{{ title }}{% endblock %}

{% block content %}
<div class="session">
  {% include "layouts/_login_header.html.tpl" %}
  <div class="main">
    <article>
      <div class="title">
        <h3>SignIn</h3>
      </div>
      <div class="content">
        <div class="sign-in-board">
          <a href={{ publicURL }}><span class="pure-button button-large button-success session-button"><span class="octicon octicon-mark-github"></span> Sign In with Github</span></a>
          <span class="message">This service does not access your private repository. If you want to management private repository, please click the button below.</span>
          <a href={{ privateURL }}><span class="pure-button button-small pure-button-primary secondary-session-button"><span class="octicon octicon-mark-github"></span> Sign In with Github Private Access</span></a>
        </div>
      </div>
    </article>
  </div>
</div>
{% endblock %}
