<script lang="ts" setup>
import { Record } from "pocketbase";

definePageMeta({
  middleware: ["not-authed-guard"],
});
const pb = usePocketbase();

const { isOpen, closeDialog, openDialog } = useDialogState();

const qrResult = ref();
const isLoading = ref(false);
const bookStatus = ref<RentedBookStatus>({
  bookExists: false,
  codeExists: false,
  isRented: false,
  parsedCode: {
    id: "",
    isbn: "",
    version: -1,
    rawCode: "",
  },
});

// This is only for key on child component.
const scannedBookId = ref("");

type MassolitCode = {
  version: number;
  id: string;
  isbn: string;
  rawCode: string;
};

function parseCode(code: string): MassolitCode {
  // VALID QR STRING FORMAT: MASSOLIT|<version int>|<iternal_code>|<isbn>

  const s = code.split("|");
  const massolitObject: MassolitCode = {
    version: 1,
    id: "",
    isbn: "",
    rawCode: code,
  };

  if (s[0] !== "MASSOLIT") {
    throw new Error("Invalid code format: invalid app name parameter");
  }

  const parsedIntVersion = parseInt(s[1]);

  if (typeof parsedIntVersion !== "number") {
    throw new Error("Invalid code format: version number not a number");
  }

  // const parsedIntISBN = parseInt(s[3]);

  // if (isNaN(parsedIntISBN)) {
  //   throw new Error("Invalid code format: version is NaN");
  // }

  switch (parsedIntVersion) {
    case 1:
      massolitObject.id = s[2];
      massolitObject.isbn = s[3];
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
  id: string;
};

type Renter = {
  id: string;
  first_name?: string;
  last_name?: string;
  middle_name?: string;
};
type BookInstance = {
  id: string;
  code: string;
  book_id: string;
};
type Rental = {
  managebac_user_id: string;
  book_instance_id: string;
};

export type RentedBookStatus = {
  codeExists: boolean;
  bookExists: boolean;
  isRented: boolean;
  parsedCode: MassolitCode;
  rental?: Rental;
  book_instance?: BookInstance;
  book?: RentedBook;
  bookId?: string;
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
  const bookRentedStatusModel: RentedBookStatus = {
    bookExists: false,
    codeExists: false,
    isRented: false,
    parsedCode,
  };
  bookRentedStatusModel.bookId = parsedCode.id;

  try {
    const rental = await findRental(parsedCode.id);

    if (rental !== undefined && rental.id) {
      bookRentedStatusModel.bookExists = true;
      bookRentedStatusModel.codeExists = true;
      bookRentedStatusModel.isRented = true;
      bookRentedStatusModel.rental = {
        managebac_user_id: rental.rented_to,
        book_instance_id: rental.book_instance,
      };
      console.log(bookRentedStatusModel);
      const bookInstance: Partial<Record & BookInstance> =
        rental.expand["book_instance"];

      bookRentedStatusModel.book_instance = {
        book_id: bookInstance.book,
        code: parsedCode.id,
        id: bookInstance.id,
      };

      const book: Record = rental.expand["book_instance"]["expand"]["book"];

      bookRentedStatusModel.book = {
        id: book.id,
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
      bookRentedStatusModel.bookExists = true;
      bookRentedStatusModel.codeExists = true;
      bookRentedStatusModel.bookId = parsedCode.id;

      // Need to tell ts that this will always be a single record.
      const book: Record = bookInstance.expand["book"];

      bookRentedStatusModel.book_instance = {
        book_id: bookInstance.book,
        code: parsedCode.id,
        id: bookInstance.id,
      };
      bookRentedStatusModel.book = {
        id: book.id,
        title: book.title,
        isbn: book.isbn,
        cover_url: book.cover_url,
      };
      throw "Bookinstance Found";
    }

    const book = await findBook(parsedCode.isbn);

    if (book !== undefined) {
      // offer to rent out as it exists.
      bookRentedStatusModel.bookExists = true;

      bookRentedStatusModel.book = {
        id: book.id,
        title: book.title,
        isbn: book.isbn,
        cover_url: book.cover_url,
      };
      bookRentedStatusModel.bookId;
      throw "Book found";
    }

    console.log("No rental, book instance nor book was found.");

    // code to get other remote info if first req came back problematic.
  } catch (err) {
    console.log(err);
  } finally {
    return bookRentedStatusModel;
  }
}

async function handleQrResult(result: string) {
  closeDialog();
  isLoading.value = true;
  qrResult.value = result;

  try {
    const c = parseCode(result);

    const bookStatusModel = await checkBookStatus(c);
    scannedBookId.value = c.id;

    bookStatus.value = bookStatusModel;

    if (bookStatusModel.isRented) {
      await navigateTo("/books/rented-by");
      return;
    }

    if (bookStatusModel.codeExists) {
      await navigateTo("/books/rent-out");
      return;
    }

    if (bookStatusModel.bookExists) {
      await navigateTo("/books/add-code");
      return;
    }

    await navigateTo("/books/add-book");

    console.log("BOOK STATUS: ", bookStatusModel);
  } catch (error) {
    console.error(error);
  } finally {
    isLoading.value = false;
  }
}

const route = useRoute();

onMounted(() => {
  const routeCode = route.params.code;
  console.log("Route code: ", routeCode);
  if (routeCode) {
    if (typeof routeCode === "string") {
      handleQrResult(routeCode);
    } else {
      handleQrResult(routeCode[0]);
    }
  }
  // handleQrResult("MASSOLIT|1|XGPHYS-1|9781447982463");
});
</script>

<template>
  <div
    class="h-screen w-screen background-gradient flex-col justify-between p-5 flex"
  >
    <div
      v-if="isLoading"
      class="w-screen h-screen fixed top-0 left-0 bg-background flex items-center justify-center"
    >
      <Icon name="line-md:loading-loop" class="w-16 h-16" />
    </div>
    <AppAlert />
    <div>
      <div class="flex items-center justify-between">
        <NuxtLink to="/">
          <img src="/images/logos/massolit-logo.png" class="w-24" alt="" />
        </NuxtLink>
        <p class="text-sm">{{ pb.authStore.model?.email }}</p>
      </div>
      <div class="mt-20">
        <div v-if="route.name === 'books'">
          <StudentBooks />
        </div>
        <NuxtPage v-else :rentedBookStatus="bookStatus" :key="scannedBookId" />
      </div>
    </div>
    <Dialog v-model:open="isOpen" class="max-h-min">
      <DialogTrigger class="sticky bottom-5" as-child>
        <Button class="w-full">Scan Book</Button>
      </DialogTrigger>
      <DialogContent class="h-full flex flex-col justify-between">
        <DialogHeader>
          <DialogTitle>Register Book</DialogTitle>
          <DialogDescription> Scan book QR code. </DialogDescription>
        </DialogHeader>
        <ScannerQrScanner @qrResult="handleQrResult" />

        <DialogFooter class="w-full flex items-center"> </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<style></style>
