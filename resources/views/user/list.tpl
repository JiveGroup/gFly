<ul>
{% for user in users %}
    <li>{{ user.ID }} - {{ user.Email }} - {{ user.Fullname }}</li>
{% endfor %}
</ul>
