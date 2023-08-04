let timeout = null;

$(document).ready(async () => {
  let checkpoints = await $.ajax("/api/checkpoints");
  let options = checkpoints.map((checkpoint) => {
    return `<option value="${checkpoint}">${checkpoint}</option>`;
  });
  $("#checkpointSelect").html(options).trigger("change");

  $("#submitParticipant").submit(async (event) => {
    event.preventDefault();

    try {
      if($("#participantID").val() == "") return;
      
      let data = {
        id: Number($("#participantID").val()),
        checkpoint: $("#checkpointSelect").val(),
        timestamp: Date.now(),
      };
      console.log(data);
      let response = await $.ajax("/api/participant", {
        method: "POST",
        data: JSON.stringify(data),
      });
      var message = response.message;
    } catch (err) {
      var message = err;
    } finally {
      $("#submitParticipant").trigger("reset");
      $("#message").html(message);

      clearTimeout(timeout);
      timeout = setTimeout(() => {
        $("#message").html("");
      }, 5000);
    }
  });
});
