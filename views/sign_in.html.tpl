{% extends "layouts/not_login.html.tpl" %}

{% block title %}{{ title }}{% endblock %}

{% block content %}
<div class="session">
  <header class="not-login">
    <div class="pure-menu pure-menu-horizontal">
      <span class="pure-menu-heading">fascia</span>
      <ul class="pure-menu-list right-align-list">
        <li class="pure-menu-item"><a href="/" class="pure-menu-link">About</a></li>
        <li class="pure-menu-item"><a href="/sign_up" class="pure-menu-link">SignUp</a></li>
        <li class="pure-menu-item"><a href="/sign_in" class="pure-menu-link">SignIn</a></li>
      </ul>
    </div>
  </header>
  <div class="main">
    <div class="sign-in-board">
      <form action="/sign_in" method="post" role="form" name="sign_in" id="sign_in" class="pure-form pure-form-stacked">
        <fieldset>
          <legend>Access your dashboard</legend>
          <input name="token" type="hidden" value="{{ token }}" />
          <div class="pure-control-group control-group">
            <input class="form-control" name="email" type="email" placeholder="email" />
          </div>

          <div class="pure-control-group control-group">
            <input class="form-control" name="password" type="password" placeholder="password" />
          </div>

          <div class="pure-controls control-group">
            Forgot your password?
          </div>
          <div class="pure-controls control-group">
            <button class="pure-button pure-button-primary session-button" type="submit">SignIn</button>
          </div>
        </fieldset>
      </form>
      <a href={{ oauthURL }}><span class="pure-button button-success session-button"><span class="octicon octicon-mark-github"></span> Sign In with Github</span></a>
      <p><a href="/sign_up"><span class="pure-button button-secondary session-button">SignUp</span></a></p>
    </div>
  </div>
</div>
{% endblock %}
