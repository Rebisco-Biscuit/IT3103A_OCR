-- CreateTable
CREATE TABLE "Student" (
    "id" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "age" INTEGER NOT NULL,
    "bio" TEXT,
    "interests" TEXT,
    "workHistory" TEXT,
    "education" TEXT,
    "certificates" TEXT DEFAULT '[]',

    CONSTRAINT "Student_pkey" PRIMARY KEY ("id")
);
