const checkCheckpoint = async () => {
  let checkpoint = $("#checkpointSelect").val();
  let url = `/api/checkpoints/${checkpoint}`;
  let response = await $.ajax(url);
  $("#lastUpdated").html(`Ažurirano: ${new Date().toLocaleString()}`);
  if (!response) {
    $("#tabela").hide();
    $("#message").html(
      "Nijedan takmičar još uvek nije prošao kroz ovu kontrolnu tačku."
    );
    return;
  }

  $("#message").html("");
  $("#tabela").show();
  console.log(response);
  let tableContent = response.map((participant) => {
    return `<tr>
      <td>${participant.id}</td>
      <td>${new Date(participant.timestamp).toLocaleString()}</td>
    </tr>`;
  });
  $("#tableBody").html(tableContent);
};

$(document).ready(async () => {
  let checkpoints = await $.ajax("/api/checkpoints");
  let options = checkpoints.map((checkpoint) => {
    return `<option value="${checkpoint}">${checkpoint}</option>`;
  });
  $("#checkpointSelect").html(options).trigger("change");

  setInterval(checkCheckpoint, 5 * 1000);
});
