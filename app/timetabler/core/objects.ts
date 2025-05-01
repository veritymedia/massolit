import { generateRandomId } from './id-gen'

// ==========================================
// Types and Constants
// ==========================================
export type AvailabilityRaw = {
  dow: number | string
  start: string
  end: string
}

export type TeacherRaw = {
  name: string
  availabilities: AvailabilityRaw[]
  subjects: string[]
}

export type ExamRaw = {
  subject: string
  start: string // format ISO 8601
  duration: string // format HH:mm or minutes
  room: string
  examCode?: string
}

export type BookedSegment = {
  teacher: Teacher
  start: Date
  end: Date
}

enum TORStatus {
  ExamComplete = 0,
  SubjectNotAllowed,
  NoAvailability,
}

// Randomness configuration
const TEACHER_SELECTION_RANDOMNESS = 0.3 // 0 = no randomness, 1 = maximum randomness

// ==========================================
// Availability Class
// ==========================================
export class Availability {
  dow: number
  start: Date
  end: Date

  constructor(avail: AvailabilityRaw) {
    this.dow = typeof avail.dow === 'number' ? avail.dow : parseInt(avail.dow)
    this.start = this.parseTime(avail.start)
    this.end = this.parseTime(avail.end)
  }

  private parseTime(timeStr: string): Date {
    const [hours, minutes] = timeStr.split(':').map(Number)
    const date = new Date(1970, 0, 1, hours, minutes) // Always use Jan 1, 1970
    return date
  }

  getFormattedTime(date: Date): string {
    return date.toTimeString().slice(0, 5) // e.g., "09:00"
  }
}

// ==========================================
// Teacher Class
// ==========================================
export class Teacher {
  id: string
  name: string
  availabilities: Availability[]
  subjects: string[]
  bias: number

  constructor(teacher: TeacherRaw) {
    this.id = generateRandomId()
    this.name = teacher.name
    this.bias = Math.floor(Math.random() * 5)
    this.subjects = teacher.subjects.map((sub: string) => {
      return sub.toLowerCase()
    })
    this.availabilities = teacher.availabilities.map((a) => {
      return new Availability(a)
    })
  }
}

// ==========================================
// WorkloadTracker Class
// ==========================================
export class WorkloadTracker {
  private teacherWorkloads = new Map<string, number>() // teacher.id -> minutes assigned

  /**
   * Add assigned time to teacher's workload
   * @param teacherId Teacher's ID
   * @param minutes Number of minutes assigned
   */
  addAssignment(teacherId: string, minutes: number): void {
    const currentWorkload = this.teacherWorkloads.get(teacherId) || 0
    this.teacherWorkloads.set(teacherId, currentWorkload + minutes)
  }

  /**
   * Get teacher's current workload in minutes
   * @param teacherId Teacher's ID
   */
  getWorkload(teacherId: string): number {
    return this.teacherWorkloads.get(teacherId) || 0
  }

  /**
   * Get teachers sorted by workload (lowest first)
   * @param teachers Array of teachers to sort
   */
  getTeachersSortedByWorkload(teachers: Teacher[]): Teacher[] {
    return [...teachers].sort((a, b) => {
      const workloadA = this.getWorkload(a.id)
      const workloadB = this.getWorkload(b.id)
      // Primary sort by workload, secondary by bias
      return workloadA !== workloadB ? workloadA - workloadB : a.bias - b.bias
    })
  }

  /**
   * Select teachers with weighted randomness, favoring those with less workload
   * @param eligibleTeachers Array of eligible teachers
   * @param randomnessFactor Amount of randomness (0-1)
   */
  selectTeacherWithRandomness(
    eligibleTeachers: Teacher[],
    randomnessFactor = TEACHER_SELECTION_RANDOMNESS,
  ): Teacher | null {
    if (!eligibleTeachers.length) return null

    // Sort by workload
    const sortedTeachers = this.getTeachersSortedByWorkload(eligibleTeachers)

    // Apply randomness - higher chance for teachers with less workload
    const totalTeachers = sortedTeachers.length
    const weights = sortedTeachers.map((_, index) => {
      // Lower index (less workload) gets higher weight
      return Math.pow(1 - index / totalTeachers, 2) + Math.random() * randomnessFactor
    })

    // Calculate cumulative weights
    const cumulativeWeights = []
    let sum = 0
    for (const weight of weights) {
      sum += weight
      cumulativeWeights.push(sum)
    }

    // Pick random value within total weight range
    const random = Math.random() * sum

    // Find corresponding teacher
    for (let i = 0; i < cumulativeWeights.length; i++) {
      if (random <= cumulativeWeights[i]) {
        return sortedTeachers[i]
      }
    }

    // Fallback to first teacher
    return sortedTeachers[0]
  }

  /**
   * Get all teachers' workload statistics
   */
  getWorkloadStats(): { teacherId: string; minutes: number }[] {
    return Array.from(this.teacherWorkloads.entries())
      .map(([teacherId, minutes]) => ({ teacherId, minutes }))
      .sort((a, b) => b.minutes - a.minutes)
  }

  /**
   * Reset all workload data
   */
  reset(): void {
    this.teacherWorkloads.clear()
  }
}

// ==========================================
// ExamManager Class
// ==========================================
export class ExamManager {
  static processedExams: Exam[] = []
  static workloadTracker = new WorkloadTracker()

  static addProcessedExam(exam: Exam): void {
    this.processedExams.push(exam)
  }

  static findOverlappingExamsInSameRoom(exam: Exam): Exam[] {
    return this.processedExams.filter(
      (processedExam) => processedExam.room === exam.room && exam.overlapsWithExam(processedExam),
    )
  }

  static resetProcessedExams(): void {
    this.processedExams = []
    this.workloadTracker.reset()
  }
}

// ==========================================
// Exam Class
// ==========================================
export class Exam {
  subject: string
  start: Date
  end: Date
  duration: number
  room: string
  bookedSegments: BookedSegment[]
  complete: boolean
  dow: number
  id: string

  constructor(exam: ExamRaw) {
    this.id = generateRandomId()
    this.subject = exam.subject.toLowerCase()
    this.start = new Date(exam.start)

    // Fix duration parsing - it could be in minutes or HH:MM format
    if (exam.duration.includes(':')) {
      const durationDate = this.parseTime(exam.duration)
      this.duration = durationDate.getHours() * 60 + durationDate.getMinutes()
    } else {
      // If just a number, treat as minutes directly
      this.duration = parseInt(exam.duration)
    }

    this.end = new Date(this.start)
    this.end.setMinutes(this.end.getMinutes() + this.duration)
    this.bookedSegments = []
    this.complete = false
    this.dow = this.start.getDay()
    this.room = exam.room
  }

  /**
   * Find suitable timeslots from available teachers
   * @param teachers Array of available teachers
   * @returns Status of the operation
   */
  findTimeslots(teachers: Teacher[]): TORStatus {
    // First, check if there are overlapping exams in the same room
    const overlappingExams = ExamManager.findOverlappingExamsInSameRoom(this)

    // If there are overlapping exams that are already complete, reuse their teachers
    if (overlappingExams.length > 0) {
      const completeOverlappingExams = overlappingExams.filter((exam) => exam.complete)

      if (completeOverlappingExams.length > 0) {
        return this.reuseTeachersFromOverlappingExams(completeOverlappingExams)
      }
    }

    // Check if the exam is already fully covered by the booked segments
    if (this.isFullyCovered()) {
      this.complete = true
      return TORStatus.ExamComplete
    }

    // Using workload-balanced teacher selection
    return this.findTimeslotsBalanced(teachers)
  }

  /**
   * Find timeslots using workload-balanced selection
   * @param teachers Array of available teachers
   * @returns Status of the operation
   */
  private findTimeslotsBalanced(teachers: Teacher[]): TORStatus {
    // Filter out teachers who teach this subject (they can't supervise)
    const eligibleTeachers = teachers.filter((teacher) => !teacher.subjects.includes(this.subject))

    // No eligible teachers found
    if (eligibleTeachers.length === 0) {
      return TORStatus.SubjectNotAllowed
    }

    // Track how many teachers we've tried
    let teachersTried = 0
    const maxAttempts = eligibleTeachers.length * 2 // Allow some re-tries

    // Continue until the exam is fully covered or we've tried all teachers
    while (!this.isFullyCovered() && teachersTried < maxAttempts) {
      // Select a teacher based on workload and randomness
      const teacher = ExamManager.workloadTracker.selectTeacherWithRandomness(eligibleTeachers)
      if (!teacher) break

      // Try to find overlap with this teacher
      const status = this.findTeacherOverlap(teacher)

      // If we successfully added this teacher's time, sort and merge segments
      if (status !== TORStatus.NoAvailability) {
        this.bookedSegments.sort((a, b) => a.start.getTime() - b.start.getTime())
        this.mergeOverlappingSegments()
      }

      teachersTried++

      // Check if the exam is now fully covered
      if (this.isFullyCovered()) {
        this.complete = true
        return TORStatus.ExamComplete
      }
    }

    // If we couldn't fully cover the exam
    return this.isFullyCovered() ? TORStatus.ExamComplete : TORStatus.NoAvailability
  }

  /**
   * Reuse teachers from overlapping exams in the same room
   * @param overlappingExams Array of overlapping exams in the same room
   * @returns Status of the operation
   */
  private reuseTeachersFromOverlappingExams(overlappingExams: Exam[]): TORStatus {
    // Sort overlapping exams by end time (latest end time first)
    overlappingExams.sort((a, b) => b.end.getTime() - a.end.getTime())

    // Try to copy segments from each overlapping exam
    for (const exam of overlappingExams) {
      // For each booked segment in the overlapping exam
      for (const segment of exam.bookedSegments) {
        // Create a new segment covering the overlap between this exam and the segment
        const latestStart = new Date(Math.max(this.start.getTime(), segment.start.getTime()))
        const earliestEnd = new Date(Math.min(this.end.getTime(), segment.end.getTime()))

        // If there's a valid overlap
        if (latestStart < earliestEnd) {
          // Add the segment if it doesn't overlap with existing segments
          if (!this.overlapsWithExistingSegments(latestStart, earliestEnd)) {
            this.bookedSegments.push({
              teacher: segment.teacher,
              start: latestStart,
              end: earliestEnd,
            })

            // Update the teacher's workload
            const minutesAssigned = (earliestEnd.getTime() - latestStart.getTime()) / (1000 * 60)
            ExamManager.workloadTracker.addAssignment(segment.teacher.id, minutesAssigned)
          }
        }
      }

      // Sort and merge segments after each exam's segments are added
      this.bookedSegments.sort((a, b) => a.start.getTime() - b.start.getTime())
      this.mergeOverlappingSegments()

      // Check if we've now covered the exam
      if (this.isFullyCovered()) {
        this.complete = true
        return TORStatus.ExamComplete
      }
    }

    // If we're here, we couldn't fully cover the exam with reused teachers
    return this.isFullyCovered() ? TORStatus.ExamComplete : TORStatus.NoAvailability
  }

  /**
   * Check if this exam overlaps with another exam
   * @param other Another exam to check against
   * @returns Whether there's an overlap
   */
  overlapsWithExam(other: Exam): boolean {
    return this.start < other.end && other.start < this.end
  }

  /**
   * Check if a particular teacher has valid availability for this exam
   * @param teacher The teacher to check
   * @returns Status of the matching process
   */
  private findTeacherOverlap(teacher: Teacher): TORStatus {
    // Teacher should not teach this subject to be eligible
    if (teacher.subjects.includes(this.subject)) {
      return TORStatus.SubjectNotAllowed
    }

    // Filter for the teacher's availability on the day of the exam
    const teacherAvail = teacher.availabilities.filter((a) => a.dow === this.dow)

    if (teacherAvail.length === 0) {
      return TORStatus.NoAvailability
    }

    // Keep track of whether we found at least one valid segment
    let foundValidSegment = false

    // Check each availability time slot for the teacher on this day
    for (const avail of teacherAvail) {
      // Get the overlapping time period between teacher availability and exam time
      const latestStart = this.getLatestStartTime(avail)
      const earliestEnd = this.getEarliestEndTime(avail)

      // If there's a valid overlap
      if (latestStart < earliestEnd) {
        // Check if this segment overlaps with any existing booked segments
        if (!this.overlapsWithExistingSegments(latestStart, earliestEnd)) {
          // Add the new segment
          this.bookedSegments.push({
            teacher: teacher,
            start: latestStart,
            end: earliestEnd,
          })

          // Update the teacher's workload
          const minutesAssigned = (earliestEnd.getTime() - latestStart.getTime()) / (1000 * 60)
          ExamManager.workloadTracker.addAssignment(teacher.id, minutesAssigned)

          foundValidSegment = true
        }
      }
    }

    // Sort booked segments by start time
    this.bookedSegments.sort((a, b) => a.start.getTime() - b.start.getTime())

    // Merge overlapping segments if any
    this.mergeOverlappingSegments()

    // Check if the exam is now completely covered
    if (this.isFullyCovered()) {
      this.complete = true
      return TORStatus.ExamComplete
    }

    return foundValidSegment ? TORStatus.ExamComplete : TORStatus.NoAvailability
  }

  /**
   * Get the later of the teacher's availability start and exam start
   */
  private getLatestStartTime(avail: Availability): Date {
    // Create dates with the same day as the exam but with the availability times
    const availStartOnExamDay = new Date(this.start)
    availStartOnExamDay.setHours(avail.start.getHours(), avail.start.getMinutes(), 0, 0)

    // Return the later of the two start times
    return new Date(Math.max(availStartOnExamDay.getTime(), this.start.getTime()))
  }

  /**
   * Get the earlier of the teacher's availability end and exam end
   */
  private getEarliestEndTime(avail: Availability): Date {
    // Create dates with the same day as the exam but with the availability times
    const availEndOnExamDay = new Date(this.start)
    availEndOnExamDay.setHours(avail.end.getHours(), avail.end.getMinutes(), 0, 0)

    // Return the earlier of the two end times
    return new Date(Math.min(availEndOnExamDay.getTime(), this.end.getTime()))
  }

  /**
   * Check if a new segment would overlap with existing booked segments
   * @param start Start time of the new segment
   * @param end End time of the new segment
   * @returns Whether there is an overlap
   */
  private overlapsWithExistingSegments(start: Date, end: Date): boolean {
    for (const segment of this.bookedSegments) {
      // Check for overlap
      if (start < segment.end && segment.start < end) {
        return true
      }
    }
    return false
  }

  /**
   * Merge any overlapping segments to simplify the coverage checking
   */
  private mergeOverlappingSegments(): void {
    if (this.bookedSegments.length <= 1) return

    const mergedSegments: BookedSegment[] = [this.bookedSegments[0]]

    for (let i = 1; i < this.bookedSegments.length; i++) {
      const current = this.bookedSegments[i]
      const lastMerged = mergedSegments[mergedSegments.length - 1]

      // If this segment overlaps with the previous one
      if (current.start <= lastMerged.end) {
        // Extend the end time of the last segment if needed
        if (current.end > lastMerged.end) {
          lastMerged.end = current.end
        }
      } else {
        // No overlap, add as a new segment
        mergedSegments.push(current)
      }
    }

    this.bookedSegments = mergedSegments
  }

  /**
   * Check if the exam is fully covered by booked segments
   * @returns Whether the exam is fully covered
   */
  private isFullyCovered(): boolean {
    // No segments = not covered
    if (this.bookedSegments.length === 0) {
      return false
    }

    // Sort segments by start time
    const sortedSegments = [...this.bookedSegments].sort(
      (a, b) => a.start.getTime() - b.start.getTime(),
    )

    // Check if the first segment starts at or before the exam start
    if (sortedSegments[0].start.getTime() > this.start.getTime()) {
      return false
    }

    // Check if the last segment ends at or after the exam end
    if (sortedSegments[sortedSegments.length - 1].end.getTime() < this.end.getTime()) {
      return false
    }

    // Check for gaps between segments
    let currentEnd = sortedSegments[0].end
    for (let i = 1; i < sortedSegments.length; i++) {
      if (sortedSegments[i].start.getTime() > currentEnd.getTime()) {
        // Gap found
        return false
      }
      currentEnd = new Date(Math.max(currentEnd.getTime(), sortedSegments[i].end.getTime()))
    }

    // If we're here, the exam is fully covered
    return true
  }

  private parseTime(timeStr: string): Date {
    const [hours, minutes] = timeStr.split(':').map(Number)
    const date = new Date(1970, 0, 1, hours, minutes) // Always use Jan 1, 1970
    return date
  }
}

// ==========================================
// Main Processing Function
// ==========================================
/**
 * Main function to process all exams and assign teachers
 * @param exams Array of exams to process
 * @param teachers Array of available teachers
 * @returns Array of processed exams
 */
export function processExams(exams: Exam[], teachers: Teacher[]): Exam[] {
  // Reset the processed exams list and workload tracker
  ExamManager.resetProcessedExams()

  // Sort exams by start time, then by room
  exams.sort((a, b) => {
    const startCompare = a.start.getTime() - b.start.getTime()
    if (startCompare === 0) {
      return a.room.localeCompare(b.room)
    }
    return startCompare
  })

  for (const exam of exams) {
    const status = exam.findTimeslots(teachers)

    // Add the exam to the processed list if it's complete
    if (exam.complete) {
      ExamManager.addProcessedExam(exam)
    }
  }

  return exams
}

// ==========================================
// Helper Functions
// ==========================================
/**
 * Get workload statistics for all teachers
 * @returns Array of teacher workload stats
 */
export function getTeacherWorkloadStats(): { teacherId: string; minutes: number }[] {
  return ExamManager.workloadTracker.getWorkloadStats()
}

/**
 * Get formatted workload report
 * @param teachers Array of teachers
 * @returns Formatted report string
 */
export function getWorkloadReport(teachers: Teacher[]): string {
  const stats = ExamManager.workloadTracker.getWorkloadStats()
  const teacherMap = new Map(teachers.map((t) => [t.id, t]))

  let report = 'Teacher Workload Report:\n'
  report += '------------------------\n'

  stats.forEach(({ teacherId, minutes }) => {
    const teacher = teacherMap.get(teacherId)
    if (teacher) {
      report += `${teacher.name}: ${minutes} minutes (${Math.round((minutes / 60) * 10) / 10} hours)\n`
    }
  })

  return report
}
