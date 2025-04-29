export function shapeData(row: any) {
  if (!row.date) return null;

  // Parse date
  const dateParts = row.date.split("/");
  if (dateParts.length !== 3) return null;
  const [month, day, year] = dateParts;

  // Parse time from the Time field
  let timeValue = "09:00"; // Default time if parsing fails

  if (row.time) {
    // Try to extract the time portion from the Time field
    const timeString = row.time.toLowerCase().trim();

    // check for level
    // This should be made more robust
    const isALevel =
      row.examCode[0] === "1" || row.examCode[0] === "4" ? false : true;
  }

  // Format date in ISO format
  const isoDate = `${year}-${month.padStart(2, "0")}-${day.padStart(2, "0")}T${row.time}:00`;

  // Parse duration
  // let formattedDuration = '01:00' // Default duration
  // if (row.duration) {
  //   let h = '00',
  //     m = '00'
  //   const durationMatch = row.duration.match(/(\d+)h\s*(\d+)m/)
  //   if (durationMatch) {
  //     h = durationMatch[1].padStart(2, '0')
  //     m = durationMatch[2].padStart(2, '0')
  //     formattedDuration = `${h}:${m}`
  //   } else if (row.duration.match(/(\d+)h/)) {
  //     h = row.duration.match(/(\d+)h/)[1].padStart(2, '0')
  //     formattedDuration = `${h}:00`
  //   } else if (row.duration.match(/(\d+)m/)) {
  //     m = row.duration.match(/(\d+)m/)[1].padStart(2, '0')
  //     formattedDuration = `00:${m}`
  //   }
  // }

  // Generate a unique ID using timestamp and index
  const uniqueId = `exam-${Date.now()}`;

  return {
    id: uniqueId,
    subject: row.subject || "",
    start: isoDate,
    duration: row.duration,
    room: "Room1",
    examCode: row["examCode"] || "",
  };
}
