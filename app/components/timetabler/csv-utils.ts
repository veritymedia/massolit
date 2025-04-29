import { generateId } from "~/utils";

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

  // Generate a unique ID using timestamp and index
  const uniqueId = generateId(10);

  return {
    id: uniqueId,
    subject: row.subject || "",
    start: isoDate,
    duration: row.duration,
    room: "Room1",
    examCode: row["examCode"] || "",
  };
}
