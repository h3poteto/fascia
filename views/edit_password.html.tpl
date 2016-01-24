{% extends "layouts/not_login.html.tpl" %}

{% block title %}{{ title }}{% endblock %}

{% block content %}
<div class="session">
  {% include "layouts/_login_header.html.tpl" %}
  <div class="main">
    <div class="sign-in-board">
      <form action="/passwords/{{ id }}/update" method="post" role="form" name="create" id="create" class="pure-form pure-form-stacked">
        <fieldset>
          <legend>Change your password</legend>
          <input name="token" type="hidden" value="{{ token }}" />
          <input name="reset-token" type="hidden" value="{{ resetToken }}" />
          <div class="pure-control-group control-group fascia-form-icon-wrapper">
            <input class="form-control" name="password" type="password" placeholder="pasword" />
            <div class="fascia-form-icon"><i class="fa fa-key" ></i></div>
          </div>
          <div class="pure-control-group control-group fascia-form-icon-wrapper">
            <input class="form-control" name="password-confirm" type="password" placeholder="pasword" />
            <div class="fascia-form-icon"><i class="fa fa-key" ></i></div>
          </div>
          <div class="pure-controls control-group">
            <button class="pure-button pure-button-primary session-button" type="submit">Change Password</button>
          </div>
        </fieldset>
      </form>
    </div>
  </div>
</div>
{% endblock %}
