{% extends "../layouts/webview.html.tpl" %}

{% block content %}
<div class="session">
  <div class="main">
    <div class="sign-in-board">
      <form action="/sign_in" method="post" role="form" name="sign_in" id="sign_in" class="pure-form pure-form-stacked">
        <fieldset>
          <input name="token" type="hidden" value="{{ token }}" />
          <input type="hidden" name="recaptcha_response" id="recaptchaResponse" />
          <div class="pure-control-group control-group fascia-form-icon-wrapper">
            <input class="form-control" name="email" type="email" placeholder="email" />
            <div class="fascia-form-icon"><i class="fa fa-user" ></i></div>
          </div>

          <div class="pure-control-group control-group fascia-form-icon-wrapper">
            <input class="form-control" name="password" type="password" placeholder="password" />
            <div class="fascia-form-icon"><i class="fa fa-key" ></i></div>
          </div>

          <div class="pure-controls control-group">
            <button class="pure-button pure-button-primary session-button" type="submit">SignIn</button>
          </div>
        </fieldset>
      </form>
      <a href={{ publicURL }}><span class="pure-button button-success session-button"><span class="octicon octicon-mark-github"></span> Sign In with Github</span></a>
      <span class="message">This service does not access your private repository. If you want to management private repository, please click the button below.</span>
      <a href={{ privateURL }}><span class="pure-button button-small pure-button-primary secondary-session-button"><span class="octicon octicon-mark-github"></span> Sign In with Github Private Access</span></a>
    </div>
  </div>
</div>
{% endblock %}
