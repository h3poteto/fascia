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
          <label for="email">Email:</label>
          <input class="form-control" name="email" type="email" />

          <label for="password">Password:</label>
          <input class="form-control" name="password" type="password" />

          <input name="token" type="hidden" value="{{ token }}" />

          <button class="pure-button pure-button-primary" type="submit">SignIn</button>
        </fieldset>
      </form>
      <a href={{ oauthURL }}>SignIn with github</a>
      <p><a href="sign_up">SignUp</a></p>
    </div>
  </div>
</div>
{% endblock %}
