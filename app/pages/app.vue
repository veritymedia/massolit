<script lang="ts" setup>
import { Record } from "pocketbase";
// const pb = usePocketbase();
definePageMeta({
  // middleware: ["not-authed-guard"],
});

const pb = usePocketbase();
// const mb = useManagebac();

const allStudents = ref();
async function searchManagebacStudents() {
  try {
    const res = await pb.send("/managebac/students?q=", {});
    allStudents.value = res.students;
    console.log(res);
  } catch (err) {
    console.log(err);
  }
}

const { isOpen, closeDialog, openDialog } = useDialogState();

const qrResult = ref();
const isLoading = ref(false);

type MassolitCode = {
  version: number;
  id: string;
  isbn: number;
};

function parseCode(code: string): MassolitCode {
  // VALID QR STRING FORMAT: MASSOLIT|<version int>|<iternal_code>|<isbn>

  const s = code.split("|");
  const massolitObject: MassolitCode = {
    version: 1,
    id: "",
    isbn: -1,
  };

  if (s[0] !== "MASSOLIT") {
    throw new Error("Invalid code format: invalid app name parameter");
  }

  const parsedIntVersion = parseInt(s[1]);

  if (typeof parsedIntVersion !== "number") {
    throw new Error("Invalid code format: version number not a number");
  }

  const parsedIntISBN = parseInt(s[3]);

  if (isNaN(parsedIntISBN)) {
    throw new Error("Invalid code format: version is NaN");
  }

  switch (parsedIntVersion) {
    case 1:
      massolitObject.id = s[2];
      massolitObject.isbn = parsedIntISBN;
      break;
    default:
      throw new Error("Invalid code format: unrecognised version value");
  }

  return massolitObject;
}

type RentedBook = {
  title: string;
  isbn: string;
  cover_url: string;
};

type Renter = {
  id: string;
  first_name?: string;
  last_name?: string;
  middle_name?: string;
};

type RentedBookStatus = {
  codeExists: boolean;
  bookExists: boolean;
  isRented: boolean;
  book?: RentedBook;
  renter?: Renter;
};

async function findRental(qrScannedCode: string): Promise<Record | undefined> {
  try {
    const rental = await pb
      .collection("rentals")
      .getFirstListItem(`book_instance.book_code="${qrScannedCode}"`, {
        expand: "book_instance,book_instance.book",
      });
    return rental;
  } catch (err) {
    return undefined;
  }
}

async function findBookInstance(
  qrScannedCode: string,
): Promise<Record | undefined> {
  try {
    const bookInstance = await pb
      .collection("book_instances")
      .getFirstListItem(`book_code="${qrScannedCode}"`, {
        expand: "book",
      });
    return bookInstance;
  } catch (err) {
    return undefined;
  }
}

async function findBook(isbn: number): Promise<Record | undefined> {
  try {
    const book = pb.collection("books").getFirstListItem(`isbn="${isbn}"`);
    return book;
  } catch (err) {
    return undefined;
  }
}

async function checkBookStatus(
  parsedCode: MassolitCode,
): Promise<RentedBookStatus> {
  // TODO: change the object to a class which dynmically computes the status based on the info present.
  const bookRentedStatus: RentedBookStatus = {
    bookExists: false,
    codeExists: false,
    isRented: false,
  };

  try {
    const rental = await findRental(parsedCode.id);

    if (rental !== undefined && rental.id) {
      bookRentedStatus.bookExists = true;
      bookRentedStatus.codeExists = true;
      bookRentedStatus.isRented = true;
      bookRentedStatus.renter = { id: rental.rented_to };
      const book: Record = rental.expand["book_instance"]["expand"]["book"];

      bookRentedStatus.book = {
        title: book.title,
        isbn: book.isbn,
        cover_url: book.cover_url,
      };

      // fetch ManageBac user.
      throw "Rental Found";
    }

    const bookInstance = await findBookInstance(parsedCode.id);

    if (bookInstance && bookInstance.id) {
      // offer to rent out as it exists.
      bookRentedStatus.bookExists = true;
      bookRentedStatus.codeExists = true;

      // Need to tell ts that this will always be a single record.
      const book: Record = bookInstance.expand["book"];

      bookRentedStatus.book = {
        title: book.title,
        isbn: book.isbn,
        cover_url: book.cover_url,
      };
      throw "Bookinstance Found";
    }

    const book = await findBook(parsedCode.isbn);

    if (book !== undefined) {
      // offer to rent out as it exists.
      bookRentedStatus.bookExists = true;

      bookRentedStatus.book = {
        title: book.title,
        isbn: book.isbn,
        cover_url: book.cover_url,
      };
      throw "Book found";
    }

    console.log("No rental, book instance nor book was found.");

    // code to get other remote info if first req came back problematic.
  } catch (err) {
    console.log(err);
  } finally {
    return bookRentedStatus;
  }
}

async function handleQrResult(result: string) {
  closeDialog();
  isLoading.value = true;
  qrResult.value = result;

  try {
    const c = parseCode(result);
    const rentedBookStatus = await checkBookStatus(c);

    console.log("BOOK STATUS: ", rentedBookStatus);
  } catch (error) {
    console.error(error);
  } finally {
    isLoading.value = false;
  }
}

onMounted(() => {
  handleQrResult("MASSOLIT|1|GPSYCH-3|12312312312");
});
</script>

<template>
  <div class="h-screen w-screen flex-col p-6 flex gap-10">
    <div
      v-if="isLoading"
      class="w-screen h-screen fixed top-0 left-0 bg-background flex items-center justify-center"
    >
      <Icon name="line-md:loading-loop" class="w-16 h-16" />
    </div>
    QR: {{ qrResult }}
    <Dialog v-model:open="isOpen" class="max-h-min">
      <DialogTrigger as-child>
        <Button class="w-full">Scan Book Code</Button>
      </DialogTrigger>
      <DialogContent class="h-full flex flex-col justify-between">
        <DialogHeader>
          <DialogTitle>Register Book</DialogTitle>
          <DialogDescription> Scan book QR code. </DialogDescription>
        </DialogHeader>
        <ScannerQrScanner @qrResult="handleQrResult" />

        <DialogFooter class="w-full flex items-center">
          <Button @click="closeDialog"> Next</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
