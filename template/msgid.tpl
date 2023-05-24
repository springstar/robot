package msg


var (

<%= for (id, name) in id2names { %>
  MSG_<%=capitalize(name)%> = uint32(<%=id%>)
<% } %>

)