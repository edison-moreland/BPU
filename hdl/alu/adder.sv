module adder #(
   parameter N = 4
) (
   input  logic [N-1:0] a, b,
   input  logic         carry_in,
   output logic [N-1:0] result,
   output logic         carry_out,
);

   assign {carry_out, result} = a + b + carry_in;
endmodule