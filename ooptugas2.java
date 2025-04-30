import java.util.ArrayList;
import java.util.Scanner;

class Data {
    String name;
    String password;
    String phone;

    public Data(String name, String password, String phone) {
        this.name = name;
        this.password = password;
        this.phone = phone;
    }
}

public class ooptugas2 {
    static ArrayList<Data> dataList = new ArrayList<>();
    static Scanner scanner = new Scanner(System.in);

    public static void inputData() {
        System.out.print("Enter Name: ");
        String name = scanner.nextLine();
        System.out.print("Enter Password: ");
        String password = scanner.nextLine();
        System.out.print("Enter Phone: ");
        String phone = scanner.nextLine();

        dataList.add(new Data(name, password, phone));
        System.out.println("New data is added.");
    }

    public static void showData() {
        if (dataList.isEmpty()) {
            System.out.println("| No data exists! |");
        } else {
            int maxBorder = 45;

            for (int i = 0; i < maxBorder; i++) {
                System.out.print("=");
            }
            System.out.println();

            System.out.printf("%-5s %-15s %-15s %-15s\n", "No", "Name", "Pass", "Phone");

            for (int i = 0; i < maxBorder; i++) {
                System.out.print("=");
            }
            System.out.println();

            for (int i = 0; i < dataList.size(); i++) {
                Data data = dataList.get(i);
                System.out.printf("%-5d %-15s %-15s %-15s\n", (i + 1), data.name, data.password, data.phone);

                for (int j = 0; j < maxBorder; j++) {
                    System.out.print("=");
                }
                System.out.println();
            }
        }
    }

    public static void deleteData() {
        if (dataList.isEmpty()) {
            System.out.println("| No data exists! |");
        } else {
            showData();
            System.out.print("Input data number to be deleted: ");
            int deleteIdx = scanner.nextInt();
            scanner.nextLine();
            if (deleteIdx > 0 && deleteIdx <= dataList.size()) {
                dataList.remove(deleteIdx - 1);
                System.out.println("Data is removed.");
            } else {
                System.out.println("Invalid number.");
            }
        }
    }

    public static void mainMenu() {
        while (true) {
            System.out.println("\n1. Input Data");
            System.out.println("2. Show Data");
            System.out.println("3. Delete Data");
            System.out.println("4. Exit");

            System.out.print("Your choice: ");
            String choice = scanner.nextLine();

            switch (choice) {
                case "1":
                    inputData();
                    break;
                case "2":
                    showData();
                    break;
                case "3":
                    deleteData();
                    break;
                case "4":
                    System.out.println("Exiting the program.");
                    return;
                default:
                    System.out.println("Invalid choice, please try again.");
            }
        }
    }

    public static void main(String[] args) {
        mainMenu();
    }
}